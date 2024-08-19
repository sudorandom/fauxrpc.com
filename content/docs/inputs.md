---
title: 'Inputs'
weight: 30
slug: inputs
---

FauxRPC offers versatile ways to define the services it emulates. Whether you have Protobuf descriptors, a live gRPC server, or even Buf Schema Registry images, you can seamlessly provide the necessary input for FauxRPC to generate mock responses. This flexibility empowers you to test and develop against various scenarios without relying on actual backend implementations.

## Using Descriptors
Make an `example.proto` file (or use a file that already exists):
```protobuf
syntax = "proto3";

package greet.v1;

message GreetRequest {
  string name = 1;
}

message GreetResponse {
  string greeting = 1;
}

service GreetService {
  rpc Greet(GreetRequest) returns (GreetResponse) {}
}
```

Create a descriptors file and use it to start the FauxRPC server:
```shell
$ buf build ./example.proto -o ./example.binpb
$ fauxrpc run --schema=./example.binpb
2024/08/17 08:01:19 INFO Listening on http://127.0.0.1:6660
2024/08/17 08:01:19 INFO See available methods: buf curl --http2-prior-knowledge http://127.0.0.1:6660 --list-methods
```
Done! It's that easy. Now you can call the service with any tooling that supports gRPC, gRPC-Web, or connect. So [buf curl](https://buf.build/docs/reference/cli/buf/curl), [grpcurl](https://github.com/fullstorydev/grpcurl), [Postman](https://www.postman.com/), [Insomnia](https://insomnia.rest/) all work fine!

```shell
$ buf curl --http2-prior-knowledge http://127.0.0.1:6660/greet.v1.GreetService/Greet
{
  "greeting":  "3 wolf moon fashion axe."
}
```

## Using Server Reflection
If there's an existing gRPC service running that you want to emulate, you can use server reflection to start the FauxRPC service:
```shell
$ fauxrpc run --schema=https://demo.connectrpc.com
```

## From BSR (Buf Schema Registry)
Buf has a [schema registry](https://buf.build/product/bsr) where many schemas are hosted. Here's how to use FauxRPC using images from the registry.

```shell
$ buf build buf.build/bufbuild/registry -o bufbuild.registry.json
$ fauxrpc run --schema=./bufbuild.registry.json
```

This will start a fake version of the BSR API by downloading descriptors for [bufbuild/registry](https://buf.build/bufbuild/registry) from the BSR and using them with FauxRPC. Very meta.

## Multiple Sources
You can define this `--schema` option as many times as you want. That means you can add services from multiple descriptors and even mix and match from descriptors and from server reflection:
```shell
$ fauxrpc run --schema=https://demo.connectrpc.com --schema=./example.binpb
```
