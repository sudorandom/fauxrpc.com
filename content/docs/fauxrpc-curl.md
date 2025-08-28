---
title: 'fauxrpc curl'
weight: 40
slug: fauxrpc-curl
---

The `fauxrpc curl` command is a powerful, built-in client for sending requests with dynamically generated fake data to any gRPC, gRPC-Web, or ConnectRPC service. It's the perfect tool for quickly testing endpoints, running smoke tests, or debugging services without needing to write a dedicated client.

It can use a schema source (like a Protobuf file or a Buf Schema Registry URL) or server reflection to automatically understand the service's API and generate valid requests.

---

## Flags

```shell
Usage: fauxrpc curl [<method>] [flags]

Make requests with fake data

Arguments:
  [<method>]    Service or method name.

Flags:
  -h, --help                            Show context-sensitive help.
  -l, --log-level="info"                Set the logging level (debug|info|warn|error)
      --version                         Print version information and quit

      --schema=SCHEMA,...               RPC schema modules.
  -a, --addr="http://127.0.0.1:6660"    Address to make a request to
  -p, --protocol="grpc"                 Protocol to use for requests.
  -e, --encoding="proto"                Encoding to use for requests.
      --http2-prior-knowledge           This flag can be used to indicate that HTTP/2 should be used.
      --http-3                          Enables HTTP/3.
```

Here is a comprehensive list of all the flags available for the `fauxrpc curl` command.

| Flag | Default | Description |
| :--- | :--- | :--- |
| **Target & Schema** |
| `[<method>]` | `(none)` | The target service or method name to call (e.g., `my.service.v1.UserService/GetUser`). If omitted, `fauxrpc curl` will attempt to call all methods in the schema. |
| `--schema=...` | `(none)` | 📜 Specifies the source for the RPC schema. It can be a local file path, a directory, or a URL. If not provided, the command will rely on server reflection. |
| `-a, --addr` | `http://127.0.0.1:6660` | The full address of the server to send the request to. |
| **Protocol & Encoding** |
| `-p, --protocol` | `grpc` | The protocol to use for the request. Valid options are `grpc`, `grpc-web`, and `connect`. |
| `-e, --encoding` | `proto` | The encoding to use for the request payload. Valid options are `proto` and `json`. |
| `--http2-prior-knowledge` | `false` | Use this flag to connect to a gRPC server that does not use TLS, indicating that the client should use HTTP/2 directly without an initial HTTP/1.1 upgrade request. |
| `--http-3` | `false` | Enables experimental HTTP/3 support for the request. |
| **General** |
| `-l, --log-level` | `info` | Sets the logging verbosity. Available levels are `debug`, `info`, `warn`, and `error`. |
| `--version` | `false` | Prints the current `fauxrpc` version information and then exits. |
| `-h, --help` | `false` | Shows a summary of the command and its flags. |

---

## Practical Examples

### 1. Test All Methods in a Service
This command uses a remote schema from the Buf Schema Registry to discover all available RPCs and sends a request with fake data to each one. This is a great way to perform a quick smoke test on an entire service.

```shell
fauxrpc curl \
  --http2-prior-knowledge \
  --schema=buf.build/bufbuild/registry
```

### Call a Specific RPC Method

To target a single method, simply provide its fully qualified name as an argument after the flags.

```bash
fauxrpc curl \
  --http2-prior-knowledge \
  --schema=buf.build/bufbuild/registry \
  buf.registry.plugin.v1beta1.LabelService/ListLabels
```

### Use Server Reflection

If the target server has reflection enabled, you don't need to provide a --schema. fauxrpc curl will automatically query the server to determine how to construct the request.

```bash
# Assuming the server at the default address has reflection enabled
fauxrpc curl \
  --http2-prior-knowledge \
  buf.registry.plugin.v1beta1.LabelService/ListLabels
```

### Target a gRPC-Web Service

You can easily switch protocols to test different server implementations, like a gRPC-Web service that uses JSON encoding.

```bash
fauxrpc curl \
  --addr=http://localhost:8080 \
  --protocol=grpc-web \
  --encoding=json \
  my.webapp.v1.AuthService/Login
```
