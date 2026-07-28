---
title: 'Run Server'
weight: 40
slug: fauxrpc-run
description: "A comprehensive guide to `fauxrpc run` and its flags for starting fake OpenAPI and Protobuf services."
icon: "play_circle"
---

The `fauxrpc run` command starts fake OpenAPI and Protobuf services for testing and development without a real backend.

## Flags

```shell
Usage: fauxrpc run [flags]

Run the FauxRPC server

Flags:
  -h, --help                     Show context-sensitive help.
  -l, --log-level="info"         Set the logging level (debug|info|warn|error)
      --version                  Print version information and quit

      --schema=SCHEMA,...        The schemas to serve. It can be protobuf descriptors (binpb, json, yaml), an OpenAPI specification, a URL, or a directory of schemas.
  -a, --addr="127.0.0.1:6660"    Address to bind to.
      --no-reflection            Disables the server reflection service.
      --no-http-log              Disables the HTTP log.
      --no-validate              Disables protovalidate.
      --no-doc-page              Disables the documentation page.
      --no-cors                  Disables CORS headers.
      --https                    Enables HTTPS, requires cert and certkey
      --cert=STRING              Path to certificate file
      --cert-key=STRING          Path to certificate key file
      --http-3                   Enables HTTP/3 support.
      --empty                    Allows the server to run with no services.
      --only-stubs               Only use pre-defined stubs and don't make up fake data.
      --stubs=STUBS,...          Directories or file paths for JSON files.
      --dashboard                Enable the admin dashboard.
      --depth=5                  Max depth for generated messages.
      --static-seed              Use deterministic generated values for unstubbed OpenAPI and Protobuf requests.
      --proxy-to=STRING          Upstream gRPC or Connect server address to proxy requests to.
      --record-dir=STRING        Directory where recorded stubs should be saved.
      --ssl-keylog-file=STRING   Path to file for logging TLS secrets; requires HTTPS or HTTP3.
```

## Flag reference

Here is a comprehensive list of all the flags available for the `fauxrpc run` command.

| Flag | Default | Description |
| :--- | :--- | :--- |
| **Schema & Data** |
| `--schema=...` | `(none)` | **(Required)** 📜 Specifies an OpenAPI or Protobuf schema source. It can be a local path, directory, or URL and may be repeated. See [Inputs](/docs/server/inputs/). |
| `--stubs=...` | `(none)` | A YAML/JSON file or directory containing predefined Protobuf or OpenAPI responses. |
| `--only-stubs` | `false` | Disables generated fallback data. Unstubbed OpenAPI operations return HTTP `501`; unstubbed Protobuf RPCs return an empty message. |
| `--empty` | `false` | Allows the server to start without any services loaded from a schema. Useful for starting a base server that might be configured dynamically. |
| `--depth` | `5` | Limits recursion depth when generating nested Protobuf messages and OpenAPI response schemas. |
| `--static-seed` | `false` | Uses stable, identity-derived seeds for unstubbed OpenAPI operations and Protobuf RPC methods. Generated values vary per request when omitted. Explicit examples, defaults, and stubs are unaffected. |
| **Proxying & Recording** |
| `--proxy-to=...` | `(none)` | 🔄 The network address of the upstream gRPC or Connect server to forward requests to. |
| `--record-dir=...` | `(none)` | 🎙️ The local directory path where recorded stubs should be saved, structured by service and method. |
| **Network & Security** |
| `-a, --addr` | `127.0.0.1:6660` | The network address and port for the server to bind to (e.g., `:8080` or `0.0.0.0:9000`). |
| `--https` | `false` | 🛡️ Enables HTTPS. Requires `--cert` and `--cert-key` to be provided. |
| `--cert` | `(none)` | The file path to your TLS certificate. Required if `--https` is used. |
| `--cert-key` | `(none)` | The file path to your TLS certificate key. Required if `--https` is used. |
| `--http-3` | `false` | Enables experimental HTTP/3 support for improved performance. Requires `--https`. |
| `--no-cors` | `false` | Disables the default Cross-Origin Resource Sharing (CORS) headers, which are permissive by default. |
| **Features & Services** |
| `--dashboard` | `false` | 📊 Enables a web-based admin dashboard for inspecting services, methods, and server state. |
| `--no-reflection` | `false` | Disables the gRPC server reflection service. This prevents clients from dynamically discovering the server's API. |
| `--no-doc-page` | `false` | Disables the built-in documentation web page that lists all available services and methods. |
| `--no-validate` | `false` | Disables request validation using `protovalidate`. This can be useful for performance or for testing how your client handles invalid data. |
| `--no-http-log` | `false` | Disables logging of incoming HTTP requests to the console, resulting in a quieter output. |
| **General** |
| `-l, --log-level` | `info` | Sets the logging verbosity. Available levels are `debug`, `info`, `warn`, and `error`. |
| `--version` | `false` | Prints the current `fauxrpc` version information and then exits. |
| `-h, --help` | `false` | Shows a summary of the command and its flags. |

## Dynamic Fakes vs. Static Stubs

FauxRPC has two primary modes for generating responses:

1.  **Dynamic Faking (Default):** If you only provide a `--schema`, FauxRPC will dynamically generate a valid, randomized response for any RPC call it receives. This is perfect for general-purpose testing where you just need *some* valid data.

2.  **Static Stubbing:** By using the `--stubs` flag, you can provide YAML or JSON files that define specific responses for RPC calls and OpenAPI operations. This is ideal for integration tests or frontend development where you need predictable and consistent data to verify application logic. `--only-stubs` disables generated fallback data for both schema types.

Use `--static-seed` for deterministic generated responses without defining stubs. The option applies to both OpenAPI operations and Protobuf RPC methods.

---

## Practical Examples

### From Protobuf Descriptor

This is the most common use case. It starts a server using a binary Protobuf descriptor file (`.binpb`).

```bash
fauxrpc run --schema=./my-service.binpb
```

### From OpenAPI

Start an HTTP mock server from an OpenAPI YAML or JSON document:

```bash
fauxrpc run --schema=./openapi.yaml
```

Interactive documentation is available at `http://127.0.0.1:6660/fauxrpc/openapi-docs/`. See [OpenAPI Support](/docs/server/openapi/) for generated responses and operation-based stubs.

### Deterministic Generated Responses

Keep generated values stable across repeated unstubbed calls:

```bash
fauxrpc run --schema=./openapi.yaml --static-seed
```

### Live Reflection

Start a server that mimics a remote gRPC server by connecting to it and using its reflection API to discover services.

```bash
fauxrpc run --schema=grpc.server.com:443
```

### Multiple Schemas

You can load services from multiple sources, such as a local file and a remote server, into a single FauxRPC instance.

```bash
fauxrpc run \
  --schema=./user-service.binpb \
  --schema=grpc.payments.com:443
```

### Predefined Stubs

Start a server that uses your schema but serves predictable data from a directory of JSON stub files for specific RPCs.

```bash
fauxrpc run \
  --schema=./inventory.binpb \
  --stubs=./testdata/stubs/
```

### Upstream Proxying

Run the server as an intercepting proxy that forwards requests to an upstream server, falling back to dynamic fake data or stubs for unimplemented endpoints:

```bash
fauxrpc run \
  --schema=buf.build/connectrpc/eliza \
  --proxy-to=127.0.0.1:8080
```

### Stub Recording

Run in proxy mode and record all passing traffic as structured JSON stubs into a local directory:

```bash
fauxrpc run \
  --schema=buf.build/connectrpc/eliza \
  --proxy-to=127.0.0.1:8080 \
  --record-dir=./testdata/stubs/
```

### Admin Dashboard

Run the server and enable the web dashboard, which is great for inspecting services and seeing traffic.
```bash
fauxrpc run \
  --schema=./my-service.binpb \
  --dashboard
```

### Secure Server (HTTPS)

To run a server on a custom port with TLS enabled, you must provide a certificate and key.

```bash
fauxrpc run \
  --addr=:8443 \
  --https \
  --cert=./certs/server.crt \
  --cert-key=./certs/server.key \
  --schema=./my-service.binpb
```

### Custom Lean Server

This example starts a minimal server with reflection and the doc page disabled, listening on all network interfaces.

```bash
fauxrpc run \
  --addr=0.0.0.0:9000 \
  --no-reflection \
  --no-doc-page \
  --schema=./my-service.binpb
```
