---
title: 'fauxrpc run'
weight: 40
slug: fauxrpc-run
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

      --schema=SCHEMA,...        The modules to use for the RPC schema. It can be protobuf descriptors (binpb, json,
                                 yaml), a URL for reflection or a directory of descriptors.
  -a, --addr="127.0.0.1:6660"    Address to bind to.
      --no-reflection            Disables the server reflection service.
      --no-doc-page              Disables the documentation page.
      --https                    Enables HTTPS, requires cert and certkey
      --cert=STRING              Path to certificate file
      --cert-key=STRING          Path to certificate key file
      --http-3                   Enables HTTP/3 support.
      --empty                    Allows the server to run with no services.
```

Let's break down each flag into a little bit more detail.

- `-h, --help`: Displays help information about the command and its flags.
- `-l, --log-level`: Sets the logging level for the server. Options are `debug`, `info`, `warn`, and `error`. Default is `info`.
- `--version`: Shows the version of `fauxrpc` you're using.
- `--schema`: Specifies the source for generating the fake gRPC services. It accepts:
    - **Protobuf descriptors:** Paths to files containing Protobuf definitions (`.binpb`, `.json`, `.yaml`). This is the most common way to define your services.
    - **Reflection URL:** A URL of a live gRPC server that supports reflection. FauxRPC will query this server to understand its services and generate fake responses accordingly.
    - **Directory:** A path to a directory containing Protobuf descriptor files.
    * You can use this flag multiple times to combine services from different sources.
    * See the [Inputs](/docs/server/inputs/) page for more on this.
- `-a, --addr`: Sets the address and port for the server to listen on. Default is `127.0.0.1:6660`.
- `--no-reflection`: Disables the server reflection service. This is useful if you don't want clients to be able to introspect the server's API.
- `--no-doc-page`: Disables the documentation page served by FauxRPC.
- `--https`: Enables HTTPS for secure communication. Requires `--cert` and `--cert-key` to be set.
- `--cert`: Specifies the path to the TLS certificate file.
- `--cert-key`: Specifies the path to the TLS certificate key file.
- `--http-3`: Enables HTTP/3 support for the server.
- `--empty`: Allows starting the server without any defined services. This can be useful for specific testing scenarios.

## Usage Examples

1. **Start a server using a Protobuf descriptor file:**

```bash
fauxrpc run --schema=./my-service.binpb
```

2. **Start a server using a live gRPC server for reflection:**

```bash
fauxrpc run --schema=https://example.com:50051 
```

3. **Start a server with services from both a descriptor and reflection:**

```bash
fauxrpc run \
  --schema=./my-service.binpb \
  --schema=https://example.com:50051
```

4. **Start a server with a specific address and HTTPS enabled:**

```bash
fauxrpc run \
  --addr=:8443 \
  --https \
  --cert=./server.crt \
  --cert-key=./server.key \
  --schema=./my-service.binpb
```

5. **Start a server with debug logging:**

```bash
fauxrpc run \
  --log-level=debug \
  --schema=./my-service.binpb
```

These examples demonstrate how to use the `fauxrpc run` command with various options to create a fake gRPC server tailored to your testing needs.
