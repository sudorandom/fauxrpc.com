---
title: FauxRPC - create a working mock gRPC server in minutes
layout: home
---
![FauxRPC](</logo-wide.jpg>)

{{< hextra/hero-headline style="font-size:1.8em" >}}Say Hello to FauxRPC{{< /hextra/hero-headline >}}

{{< hextra/hero-subtitle >}}
FauxRPC is a powerful tool that empowers you to accelerate development and testing by effortlessly generating fake implementations of gRPC, gRPC-Web, Connect, and REST services. If you have a protobuf-based workflow, this tool could help.
{{< /hextra/hero-subtitle >}}

{{< hextra/hero-button text="Documentation" link="docs/intro/" >}}
{{< hextra/hero-button text="Install" link="docs/install/" >}}

{{< hextra/hero-headline style="font-size:1.3em" >}}Effortlessly start a fake service{{< /hextra/hero-headline >}}
{{< hextra/hero-container cols=1 >}}
{{< p >}}
Define your protobuf file and start FauxRPC. Now you instantly have a working service to start integrating against! You can use tools like **[grpcurl](https://github.com/fullstorydev/grpcurl)**, **[buf curl](https://buf.build/docs/reference/cli/buf/curl/)** or even **[curl](https://curl.se/)**.
{{< /p >}}
{{< tabs items="greet.proto,Run FauxRPC,Use" >}}
{{< tab >}}
This is an example protobuf file.
```protobuf
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
```
{{< /tab >}}
{{< tab >}}
```shell
$ fauxrpc run --schema greet.proto
FauxRPC (0.2.0 (697fec09ce17947605f1014095691bba43dac2ce) @ 2024-11-16T14:40:32Z; go1.23.3) - 1 services loaded
Listening on http://127.0.0.1:6660
OpenAPI documentation: http://127.0.0.1:6660/fauxrpc/openapi.html

Example Commands:
$ buf curl --http2-prior-knowledge http://127.0.0.1:6660 --list-methods
$ buf curl --http2-prior-knowledge http://127.0.0.1:6660/[METHOD_NAME]
Server started.
```
{{< /tab >}}
{{< tab >}}
Now that you have FauxRPC running, you can make requests to it! Let's start with grpcurl:
```shell
$ grpcurl -plaintext 127.0.0.1:6660 list
examples.basic.Greeter

$ grpcurl -plaintext 127.0.0.1:6660 list examples.basic.Greeter
examples.basic.Greeter.SayHello
examples.basic.Greeter.WriteHello

$ grpcurl -d '{"name": "Kevin"}' -plaintext 127.0.0.1:6660 examples.basic.Greeter/SayHello
{
  "message": "Poutine."
}
```

You can also use **buf curl**:
```shell
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
```
And obviously, generated gRPC/gRPC-Web/Connect clients all work as expected as well.
{{< /tab >}}
{{< /tabs >}}
{{< /hextra/hero-container >}}

{{< hextra/hero-headline style="font-size:1.3em" >}}REST{{< /hextra/hero-headline >}}
{{< hextra/hero-container cols=1 >}}
{{< p >}}
You can use **[google.api.http](https://github.com/googleapis/googleapis/blob/master/google/api/http.proto)** annotations to define REST-like behavior for protobuf methods. FauxRPC can use these as well to expose your REST APIs.
{{< /p >}}
{{< tabs items="users.proto,buf.yaml,Run FauxRPC,Use" >}}
{{< tab >}}
This is an example protobuf file that includes an annotation, mapping `examples.users.UserService/GetUser` to `GET /user/{user_id}`.
```protobuf
syntax = "proto3";
package examples.users;

import "google/api/annotations.proto";
import "google/protobuf/empty.proto";

service UserService {
  rpc GetUser(UserID) returns (User) {
    option (google.api.http).get = "/user/{user_id}";
    option idempotency_level = NO_SIDE_EFFECTS;
  }
}
```
{{< /tab >}}
{{< tab >}}
Because we have a remote dependency, we use buf to manage this for us `buf.yaml`:
```yaml
version: v2
deps:
 - buf.build/googleapis/googleapis
```
{{< /tab >}}
{{< tab >}}
Run the FauxRPC server. Instead of passing the protobuf file(s) as `--schema` we are building the **[protobuf descriptors](https://buf.build/docs/reference/descriptors/)** using `buf build`. This includes all remote repositories from the **[buf schema registry](https://buf.build/product/bsr)**.
```shell
$ buf dep update
$ buf build . -o users.binpb
$ fauxrpc run --schema users.binpb
FauxRPC (0.2.0 (697fec09ce17947605f1014095691bba43dac2ce) @ 2024-11-16T14:40:32Z; go1.23.3) - 1 services loaded
Listening on http://127.0.0.1:6660
OpenAPI documentation: http://127.0.0.1:6660/fauxrpc/openapi.html

Example Commands:
$ buf curl --http2-prior-knowledge http://127.0.0.1:6660 --list-methods
$ buf curl --http2-prior-knowledge http://127.0.0.1:6660/[METHOD_NAME]
Server started.
```
{{< /tab >}}
{{< tab >}}
Now that you have FauxRPC running, you can make requests to it! Let's start with grpcurl:
```shell
$ curl http://127.0.0.1:6660/user/1234
{"id":"7430821210685909574", "name":"Davin Rogahn", "status":"Selfies."}
```

You can still use gRPC/gRPC-Web or Connect:
```shell
$ buf curl -d '{}' --http2-prior-knowledge http://127.0.0.1:6660/examples.users.UserService/GetUser
{
  "id": "5798700489997268478",
  "name": "Cyril Rempel",
  "status": "Post-ironic."
}
```
And obviously, generated gRPC/gRPC-Web/Connect clients all work as expected as well.
{{< /tab >}}
{{< /tabs >}}
{{< /hextra/hero-container >}}

{{< hextra/hero-headline style="font-size:1.3em" >}}Protovalidate{{< /hextra/hero-headline >}}
{{< hextra/hero-container cols=1 >}}
{{< p >}}
You can use **[protovalidate](https://buf.build/bufbuild/protovalidate/docs/main:buf.validate)** annotations to define constraints for protobuf fields and messages. FauxRPC can use these to both produce more realistic data and for validation of requests.
{{< /p >}}
{{< tabs items="users.proto,buf.yaml,Run FauxRPC,Use" >}}
{{< tab >}}
This is an example protobuf file that includes a few protovalidate annotations.
```protobuf
syntax = "proto3";
package examples.protovalidate;

import "buf/validate/validate.proto";

message GreetRequest {
  string name = 1 [(buf.validate.field).string = {
    min_len: 3
    max_len: 100
  }];
}

message GreetResponse {
  string greeting = 1 [(buf.validate.field).string.example = "Hello, user!"];
}

service GreetService {
  rpc Greet(GreetRequest) returns (GreetResponse) {}
}
```
{{< /tab >}}
{{< tab >}}
Because we have a remote dependency, we use buf to manage this for us `buf.yaml`:
```yaml
version: v2
deps:
 - buf.build/bufbuild/protovalidate
```
{{< /tab >}}
{{< tab >}}
Run the FauxRPC server. Instead of passing the protobuf file(s) as `--schema` we are building the **[protobuf descriptors](https://buf.build/docs/reference/descriptors/)** using `buf build`. This includes all remote repositories from the **[buf schema registry](https://buf.build/product/bsr)**.
```shell
$ buf dep update
$ buf build . -o protovalidate.binpb
$ fauxrpc run --schema protovalidate.binpb
FauxRPC (0.2.0 (697fec09ce17947605f1014095691bba43dac2ce) @ 2024-11-16T14:40:32Z; go1.23.3) - 1 services loaded
Listening on http://127.0.0.1:6660
OpenAPI documentation: http://127.0.0.1:6660/fauxrpc/openapi.html

Example Commands:
$ buf curl --http2-prior-knowledge http://127.0.0.1:6660 --list-methods
$ buf curl --http2-prior-knowledge http://127.0.0.1:6660/[METHOD_NAME]
Server started.
```
{{< /tab >}}
{{< tab >}}
```shell
$ buf curl -d '{"name": "Kevin"}' --http2-prior-knowledge http://127.0.0.1:6660/examples.protovalidate.GreetService/Greet
{
  "greeting": "Hello, user!"
}
```

If you don't give a value for "name" you will receive a validation error:
```shell
$ buf curl -d '{}' --http2-prior-knowledge http://127.0.0.1:6660/examples.protovalidate.GreetService/Greet
{
   "code": "invalid_argument",
   "message": "validation error:\n - name: value length must be at least 3 characters [string.min_len]",
   "details": [
      {
         "type": "buf.validate.Violations",
         "value": "CkIKBG5hbWUSDnN0cmluZy5taW5fbGVuGip2YWx1ZSBsZW5ndGggbXVzdCBiZSBhdCBsZWFzdCAzIGNoYXJhY3RlcnM",
         "debug": {
            "violations": [
               {
                  "fieldPath": "name",
                  "constraintId": "string.min_len",
                  "message": "value length must be at least 3 characters"
               }
            ]
         }
      }
   ]
}
```
{{< /tab >}}
{{< /tabs >}}
{{< /hextra/hero-container >}}

{{< hextra/feature-grid >}}
  {{< hextra/feature-card
    title="Faster Development & Testing"
    link="docs/intro/"
    subtitle="Work independently without relying on fully functional backend services."
    icon="fast-forward"
  >}}
  {{< hextra/feature-card
    title="Isolation & Control"
    link="docs/workflow/"
    subtitle="Test frontend components in isolation with controlled fake data."
    icon="cog"
  >}}
  {{< hextra/feature-card
    title="Multi-Protocol Support"
    link="/docs/server/multi-protocol-support/"
    subtitle="Supports multiple protocols: gRPC, gRPC-Web, Connect, and REST."
    icon="globe-alt"
  >}}
  {{< hextra/feature-card
    title="Prototyping & Demos"
    link="docs/workflow/"
    subtitle="Create prototypes and demos quickly without building the full backend. Fake it till you make it."
    icon="device-mobile"
  >}}
  {{< hextra/feature-card
    title="Testcontainers support"
    link="docs/testcontainers/"
    subtitle="(go only) Superpower testing Go services by using FauxRPC with Testcontainers."
    icon="check"
  >}}
  {{< hextra/feature-card
    title="Plays well with others"
    link="/docs/protovalidate/"
    subtitle="Test data from FauxRPC will try to automatically follow any **protovalidate** constraints that are defined."
    icon="sparkles"
  >}}
{{< /hextra/feature-grid >}}

{{< hextra/hero-button text="Get Started" link="docs/intro/" >}}
