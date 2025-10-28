---
title: FauxRPC - create a working mock gRPC server in minutes
layout: home
---
{{< fauxrpc_logo >}}

## Say Hello to FauxRPC

> FauxRPC is a powerful tool that empowers you to accelerate development and testing by effortlessly generating fake implementations of gRPC, gRPC-Web, Connect, and REST services. If you have a protobuf-based workflow, this tool could help.

[Documentation](docs/intro/) | [Install](docs/install/)

### Effortlessly start a fake service

Define your protobuf file and start FauxRPC. Now you instantly have a working service to start integrating against! You can use tools like **[grpcurl](https://github.com/fullstorydev/grpcurl)**, **[buf curl](https://buf.build/docs/reference/cli/buf/curl/)** or even **[curl](https://curl.se/)**.

<div class="responsive-code-blocks">
  <div class="code-block-item">
    <h3>greet.proto</h3>
    {{< highlight protobuf >}}
syntax = "proto3";

package examples.basic;

service Greeter {
  rpc SayHello(HelloRequest) returns (HelloReply) {
    option idempotency_level = NO_SIDE_EFFECTS;
  }
  rpc WriteHello(HelloRequest) returns (HelloReply) {}
}

message HelloRequest {
  string name = 1;
  uint32 hello_count = 2;
}

message HelloReply {
  string message = 1;
}
{{< /highlight >}}
  </div>
  <div class="code-block-item">
    <h3>Run FauxRPC</h3>
    {{< highlight shell >}}
$ fauxrpc run --schema greet.proto
FauxRPC (0.2.0 (697fec09ce17947605f1014095691bba43dac2ce) @ 2024-11-16T14:40:32Z; go1.23.3) - 1 services loaded
Listening on http://127.0.0.1:6660
OpenAPI documentation: http://127.0.0.1:6660/fauxrpc/openapi.html

Example Commands:
$ buf curl --http2-prior-knowledge http://127.0.0.1:6660 --list-methods
$ buf curl --http2-prior-knowledge http://127.0.0.1:6660/[METHOD_NAME]
Server started.
{{< /highlight >}}
  </div>
  <div class="code-block-item">
    <h3>Use</h3>
    <p>Now that you have FauxRPC running, you can make requests to it! Let's start with grpcurl:</p>
    {{< highlight shell >}}
$ grpcurl -plaintext 127.0.0.1:6660 list
examples.basic.Greeter

$ grpcurl -plaintext 127.0.0.1:6660 list examples.basic.Greeter
examples.basic.Greeter.SayHello
examples.basic.Greeter.WriteHello

$ grpcurl -d '{"name": "Kevin"}' -plaintext 127.0.0.1:6660 examples.basic.Greeter/SayHello
{
  "message": "Poutine."
}
{{< /highlight >}}
    <p>You can also use <strong>buf curl</strong>:</p>
    {{< highlight shell >}}
$ buf curl --http2-prior-knowledge http://127.0.0.1:6660 --list-methods
examples.basic.Greeter/SayHello
examples.basic.Greeter/WriteHello

$ buf curl -d '{"name": "Kevin"}' --http2-prior-knowledge http://127.0.0.1:6660/examples.basic.Greeter/SayHello
{
  "message": "Asymmetrical skateboard."
}

$ buf curl -d '{"name": "Kevin"}' --protocol=grpc --http2-prior-knowledge http://127.0.0.1:6660/examples.basic.Greeter/SayHello
{
  "message": "Asymmetrical skateboard."
}

$ buf curl -d '{"name": "Kevin"}' --protocol=grpcweb --http2-prior-knowledge http://127.0.0.1:6660/examples.basic.Greeter/SayHello
{
  "message": "Asymmetrical skateboard."
}
{{< /highlight >}}
    <p>And obviously, generated gRPC/gRPC-Web/Connect clients all work as expected as well.</p>
  </div>
</div>

*   [**Faster Development & Testing**](docs/intro/): Work independently without relying on fully functional backend services.
*   [**Interactive Dashboard**](docs/server/dashboard/): Gain real-time insights into your server's operations.
*   [**Multi-Protocol Support**](/docs/server/multi-protocol-support/): Supports multiple protocols: gRPC, gRPC-Web, Connect, and REST.
*   [**Prototyping & Demos**](docs/workflow/): Create prototypes and demos quickly without building the full backend. Fake it till you make it.
*   [**Testcontainers support**](docs/testcontainers/): (go only) Superpower testing Go services by using FauxRPC with Testcontainers.
*   [**Plays well with others**](/docs/protovalidate/): Test data from FauxRPC will try to automatically follow any **protovalidate** constraints that are defined.

[Get Started](docs/intro/)
