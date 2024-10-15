---
title: 'Multi-protocol Support'
weight: 15
slug: multi-protocol-support
aliases:
- /docs/multi-protocol-support/
---

## gRPC/gRPC-Web/Connect
The multi-protocol support [is based on ConnectRPC](https://connectrpc.com/docs/multi-protocol/). So with FauxRPC, you get **gRPC, gRPC-Web and Connect** out of the box.

## REST
FauxRPC does a little bit more than gRPC-based protocols. It allows you to use [`google.api.http` annotations](https://grpc-ecosystem.github.io/grpc-gateway/docs/tutorials/adding_annotations/) to present a JSON/HTTP API, so you can gRPC and REST together! This is normally done with [an additional service](https://github.com/grpc-ecosystem/grpc-gateway) that runs in-between the outside world and your actual gRPC service but with FauxRPC you get the so-called transcoding from HTTP/JSON to gRPC all in the same package. Here's a concrete example:

```protobuf
syntax = "proto3";

package http.service;

import "google/api/annotations.proto";

service HTTPService {
  rpc GetMessage(GetMessageRequest) returns (Message) {
    option (google.api.http) = {get: "/v1/{name=messages/*}"};
  }
}
message GetMessageRequest {
  string name = 1; // Mapped to URL path.
}
message Message {
  string text = 1; // The resource content.
}
```

Again, we start the service by building the descriptors and using
```
$ buf build ./httpservice.proto -o ./httpservice.binpb
$ fauxrpc run --schema=httpservice.binpb
```

Now that we have the server running we can test this with the "normal" curl:
```shell
$ curl http://127.0.0.1:6660/v1/messages/123456
{"text":"Retro."}⏎
```
Sweet. You can now easily support REST alongside gRPC. If you are wondering how to do this with "real" services, look into [vanguard-go](https://github.com/connectrpc/vanguard-go). This library is doing the real heavy lifting.

More documentation on the `google.http.api` annotations can be found on [Google's Transcoding Documentation](https://cloud.google.com/endpoints/docs/grpc/transcoding), the [google/api/http.proto file](https://github.com/googleapis/googleapis/blob/master/google/api/http.proto) and the [gRPC-Gateway documentation site](https://grpc-ecosystem.github.io/grpc-gateway/docs/mapping/examples/).
