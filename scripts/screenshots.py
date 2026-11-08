#!/usr/bin/env -S uv run --script
# /// script
# requires-python = ">=3.11"
# dependencies = [
#   "playwright>=1.49.0",
# ]
# ///
"""Capture FauxRPC dashboard screenshots into assets/images/."""

from __future__ import annotations

import shutil
import socket
import subprocess
import tempfile
import time
import urllib.error
import urllib.request
from pathlib import Path
from playwright.sync_api import sync_playwright

ROOT = Path(__file__).resolve().parent.parent
IMAGES_DIR = ROOT / "assets" / "images"

STUBS_YAML = """\
stubs:
  - id: item-featured-keyboard
    target: examples.demo.DemoService/GetItem
    active_if: 'req.id == "550e8400-e29b-41d4-a716-446655440000"'
    priority: 100
    cel_content: |
      {
        "name": "Split Mechanical Keyboard Pro",
        "quantity": 42,
        "type": 1,
        "tags": ["hardware", "ergonomic", "hot_swap"],
        "attributes": {"switch": "tactile_v2", "layout": "columnar_stagger"},
        "status_message": "in_stock_ready_to_ship"
      }

  - id: item-digital-license
    target: examples.demo.DemoService/GetItem
    active_if: 'req.id.startsWith("00000000-")'
    priority: 80
    content:
      name: "FauxRPC Enterprise License"
      quantity: 250
      type: ITEM_TYPE_DIGITAL
      tags: ["license", "annual"]
      attributes:
        tier: "unlimited"
      status_code: 200

  - id: item-rate-limited
    target: examples.demo.DemoService/GetItem
    active_if: 'req.id == "ffffffff-ffff-ffff-ffff-ffffffffffff"'
    priority: 90
    error_code: 8
    error_message: "rate limit exceeded for item lookup"

  - id: greet-vip-welcome
    target: examples.protovalidate.GreetService/Greet
    active_if: 'req.name == "Ada Lovelace"'
    priority: 50
    cel_content: '{"greeting": "Welcome back, " + req.name + "! All 4 protocols are online."}'

  - id: user-admin-profile
    target: examples.rest.UserService/GetUser
    active_if: "req.user_id == 42"
    priority: 60
    content:
      id: 42
      name: "Ada Lovelace"
      status: "Staff Systems Architect"
"""


def wait_for_server(url: str, timeout: float = 15.0) -> None:
    deadline = time.time() + timeout
    while time.time() < deadline:
        try:
            with urllib.request.urlopen(url, timeout=1.0) as resp:
                if resp.status < 500:
                    return
        except Exception:
            time.sleep(0.15)
    raise RuntimeError(f"Timed out waiting for FauxRPC server at {url}")


def find_port(preferred: int = 6660) -> int:
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        if s.connect_ex(("127.0.0.1", preferred)) != 0:
            return preferred
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.bind(("127.0.0.1", 0))
        return s.getsockname()[1]


def http_request(
    url: str,
    method: str = "GET",
    data: bytes | None = None,
    headers: dict[str, str] | None = None,
) -> None:
    req = urllib.request.Request(url, data=data, method=method, headers=headers or {})
    try:
        with urllib.request.urlopen(req, timeout=5.0) as resp:
            resp.read()
    except urllib.error.HTTPError:
        pass


def send_sample_traffic(addr: str) -> None:
    base_url = f"http://{addr}"

    # 1. gRPC call via fauxrpc curl
    subprocess.run(
        [
            "fauxrpc",
            "curl",
            f"--addr={base_url}",
            "--http2-prior-knowledge",
            "--protocol=grpc",
            "examples.basic.Greeter/SayHello",
        ],
        check=True,
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL,
    )
    time.sleep(0.08)

    # 2. HTTP / REST call (transcoded google.api.http)
    http_request(f"{base_url}/user/42")
    time.sleep(0.08)

    # 3. ConnectRPC validation error (short name violates min_len: 3)
    http_request(
        f"{base_url}/examples.protovalidate.GreetService/Greet",
        method="POST",
        data=b'{"name": "Al"}',
        headers={
            "Content-Type": "application/json",
            "Connect-Protocol-Version": "1",
        },
    )
    time.sleep(0.08)

    # 4. gRPC-Web call via fauxrpc curl
    subprocess.run(
        [
            "fauxrpc",
            "curl",
            f"--addr={base_url}",
            "--http2-prior-knowledge",
            "--protocol=grpcweb",
            "examples.demo.DemoService/GetContactProfile",
        ],
        check=True,
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL,
    )
    time.sleep(0.08)

    # 5. ConnectRPC call matching a CEL stub
    http_request(
        f"{base_url}/examples.protovalidate.GreetService/Greet",
        method="POST",
        data=b'{"name": "Ada Lovelace"}',
        headers={
            "Content-Type": "application/json",
            "Connect-Protocol-Version": "1",
        },
    )
    time.sleep(0.08)

    # 6. HTTP REST call to /v1/items/{id}
    http_request(f"{base_url}/v1/items/550e8400-e29b-41d4-a716-446655440000")
    time.sleep(0.08)

    # 7. gRPC call to DemoService/GetContactProfile (latest top row in the stream)
    subprocess.run(
        [
            "fauxrpc",
            "curl",
            f"--addr={base_url}",
            "--http2-prior-knowledge",
            "--protocol=grpc",
            "examples.demo.DemoService/GetContactProfile",
        ],
        check=True,
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL,
    )


def main() -> None:
    if not shutil.which("fauxrpc"):
        raise SystemExit("fauxrpc binary not found on PATH. Run `mise install` first.")
    if not shutil.which("buf"):
        raise SystemExit("buf binary not found on PATH. Run `mise install` first.")

    IMAGES_DIR.mkdir(parents=True, exist_ok=True)
    port = find_port(6660)
    addr = f"127.0.0.1:{port}"
    dashboard_url = f"http://{addr}/fauxrpc"

    with tempfile.TemporaryDirectory(prefix="fauxrpc-screenshots-") as tmp:
        tmp_dir = Path(tmp)
        schema_binpb = tmp_dir / "schema.binpb"
        stubs_file = tmp_dir / "stubs.yaml"
        stubs_file.write_text(STUBS_YAML, encoding="utf-8")

        print("Building example Protobuf descriptor set...")
        subprocess.run(
            ["buf", "build", "-o", str(schema_binpb)],
            cwd=ROOT,
            check=True,
        )

        print(f"Starting FauxRPC server on {addr}...")
        server = subprocess.Popen(
            [
                "fauxrpc",
                "run",
                f"--schema={schema_binpb}",
                f"--stubs={stubs_file}",
                "--dashboard",
                f"--addr={addr}",
            ],
            cwd=ROOT,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            text=True,
        )

        try:
            wait_for_server(dashboard_url)

            with sync_playwright() as p:
                browser = p.chromium.launch(headless=True)
                context = browser.new_context(
                    viewport={"width": 1440, "height": 840},
                    device_scale_factor=2,
                    color_scheme="dark",
                )
                page = context.new_page()

                # Open Request Log first so SSE stream captures live traffic
                page.goto(f"{dashboard_url}/", wait_until="networkidle")
                page.wait_for_timeout(400)

                print("Sending multi-protocol sample traffic...")
                send_sample_traffic(addr)
                page.wait_for_timeout(700)

                # Select the top request row so the inspector pane is populated
                first_row = page.locator("tr.log-entry-row").first
                first_row.wait_for(state="visible", timeout=5000)
                first_row.click()
                page.wait_for_timeout(300)

                # 1. Request Log screenshot (Master-detail request stream + payload inspector)
                print("Capturing dashboard-request-log.png and dashboard-request-log-home.png...")
                page.screenshot(path=str(IMAGES_DIR / "dashboard-request-log.png"))
                page.screenshot(path=str(IMAGES_DIR / "dashboard-request-log-home.png"))

                # 2. Dashboard overview screenshot (with live server telemetry popover open)
                print("Capturing dashboard.png...")
                status_btn = page.locator("#status-popover-btn")
                if status_btn.count() > 0:
                    status_btn.click()
                    page.wait_for_timeout(250)
                page.screenshot(path=str(IMAGES_DIR / "dashboard.png"))

                # Close status popover if open
                if status_btn.count() > 0:
                    status_btn.click()
                    page.wait_for_timeout(150)

                # 3. Schema Browser screenshots
                print("Capturing dashboard-schema1.png and dashboard-schema2.png...")
                page.goto(f"{dashboard_url}/browser/", wait_until="networkidle")
                page.wait_for_timeout(350)
                demo_row = page.locator(
                    'button.schema-file-row[data-schema-search*="demo.proto"]'
                ).first
                if demo_row.count() > 0:
                    demo_row.click()
                    page.wait_for_timeout(300)
                page.screenshot(path=str(IMAGES_DIR / "dashboard-schema1.png"))

                users_row = page.locator(
                    'button.schema-file-row[data-schema-search*="users.proto"]'
                ).first
                if users_row.count() > 0:
                    users_row.click()
                    page.wait_for_timeout(300)
                page.screenshot(path=str(IMAGES_DIR / "dashboard-schema2.png"))

                # 4. Stubs Inspector screenshot
                print("Capturing dashboard-stubs.png...")
                page.goto(f"{dashboard_url}/stubs/", wait_until="networkidle")
                page.wait_for_timeout(350)
                stub_row = page.locator("tr.stub-entry-row").first
                if stub_row.count() > 0:
                    stub_row.click()
                    page.wait_for_timeout(300)
                page.screenshot(path=str(IMAGES_DIR / "dashboard-stubs.png"))

                browser.close()
                print("All dashboard screenshots updated in assets/images/.")
        finally:
            server.terminate()
            try:
                server.wait(timeout=5)
            except subprocess.TimeoutExpired:
                server.kill()


if __name__ == "__main__":
    main()
