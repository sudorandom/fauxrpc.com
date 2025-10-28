---
title: 'fauxrpc run'
weight: 40
slug: fauxrpc-run
description: "A comprehensive guide to the `fauxrpc run` command and all its flags for starting a fake gRPC server."
icon: "play_circle"
---

The `fauxrpc run` command starts a fake gRPC server that you can use for testing and development without a real backend. Here's a breakdown of its options:

## Flags

```shell
Usage: fauxrpc run [flags]

Run the FauxRPC server

Flags:
  -h, --help                     Show context-sensitive help.
  -l, --log-level="info"         Set the logging level (debug|info|warn|error)
      --version                  Print version information and quit

      --schema=SCHEMA,...        The modules to use for the RPC schema. It can be protobuf descriptors (binpb, json, yaml), a URL for reflection or a directory of descriptors.
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
```

## Flags

Here is a comprehensive list of all the flags available for the `fauxrpc run` command.

| Flag | Default | Description |
| :--- | :--- | :--- |
| **Schema & Data** |
| `--schema=...` | `(none)` | **(Required)** 📜 Specifies the source for the RPC schema. It can be a local file path, a directory, or a URL. You can use this flag multiple times. See the [Inputs](/docs/server/inputs/) page for details. |
| `--stubs=...` | `(none)` | A path to a directory or specific JSON files containing predefined responses. Use this to serve specific, static data for certain RPC calls. |
| `--only-stubs` | `false` | If enabled, the server will **only** respond with data from loaded stubs. RPC calls without a matching stub will return an error instead of dynamically generated fake data. |
| `--empty` | `false` | Allows the server to start without any services loaded from a schema. Useful for starting a base server that might be configured dynamically. |
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

2.  **Static Stubbing:** By using the `--stubs` flag, you can provide JSON files that define specific responses for specific RPC calls. This is ideal for integration tests or frontend development where you need predictable and consistent data to verify application logic. Using `--only-stubs` ensures that *only* your predefined stubs are used.

---

## Practical Examples

### Basic Server from a Protobuf Descriptor
This is the most common use case. It starts a server using a binary Protobuf descriptor file (`.binpb`).

```bash
fauxrpc run --schema=./my-service.binpb
```

### Mirroring a Live Server via Reflection

Start a server that mimics a remote gRPC server by connecting to it and using its reflection API to discover services.

```bash
fauxrpc run --schema=grpc.server.com:443
```

### Combining Multiple Schema Sources

You can load services from multiple sources, such as a local file and a remote server, into a single FauxRPC instance.

```bash
fauxrpc run \
  --schema=./user-service.binpb \
  --schema=grpc.payments.com:443
```

### Serving Predefined Responses from Stubs

Start a server that uses your schema but serves predictable data from a directory of JSON stub files for specific RPCs.

```bash
fauxrpc run \
  --schema=./inventory.binpb \
  --stubs=./testdata/stubs/
```

### Enabling the Admin Dashboard

Run the server and enable the web dashboard, which is great for inspecting services and seeing traffic.
```bash
fauxrpc run \
  --schema=./my-service.binpb \
  --dashboard
```

### Running a Secure Server with HTTPS

To run a server on a custom port with TLS enabled, you must provide a certificate and key.

```bash
fauxrpc run \
  --addr=:8443 \
  --https \
  --cert=./certs/server.crt \
  --cert-key=./certs/server.key \
  --schema=./my-service.binpb
```

### Running a Lean Server with a Custom Address

This example starts a minimal server with reflection and the doc page disabled, listening on all network interfaces.

```bash
fauxrpc run \
  --addr=0.0.0.0:9000 \
  --no-reflection \
  --no-doc-page \
  --schema=./my-service.binpb
```
