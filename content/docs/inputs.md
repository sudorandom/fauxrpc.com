---
title: 'Inputs'
weight: 30
slug: inputs
aliases:
- /docs/inputs/
- /docs/server/inputs/
description: "Provide OpenAPI documents, Protobuf files, descriptor sets, server reflection endpoints, and Buf Schema Registry modules to FauxRPC."
icon: "input"
---

FauxRPC uses one repeatable `--schema` option for OpenAPI documents, Protobuf descriptor sets, live gRPC reflection endpoints, and Buf Schema Registry modules. Different source types can be served together.

## From OpenAPI documents

OpenAPI 3 specifications can be loaded from YAML or JSON files:

```shell
$ fauxrpc run --schema=./openapi.yaml
```

HTTP(S) URLs and directories containing OpenAPI documents are also supported. FauxRPC detects the document type from its contents, serves its HTTP operations, and exposes interactive documentation at `/fauxrpc/openapi-docs/`.

See [OpenAPI Support](/docs/server/openapi/) for response generation, request validation, stubs, and `--static-seed` behavior.

## From Protobuf files
If you have `.proto` source files, you can use `buf` to build a descriptor set and then use it with FauxRPC.

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

## Using Descriptor Sets
If you already have a Protobuf descriptor set file (commonly with a `.binpb` or `.json` extension), you can provide it directly to FauxRPC. For more about descriptor sets, see the [Buf documentation](https://buf.build/docs/reference/descriptors).

```shell
$ fauxrpc run --schema=./example.binpb
```

## Using Server Reflection
If there's an existing gRPC service running that you want to emulate, you can use server reflection to start the FauxRPC service:
```shell
$ fauxrpc run --schema=https://demo.connectrpc.com
```

## From BSR (Buf Schema Registry)
Buf has a [schema registry](https://buf.build/product/bsr) where many schemas are hosted. You can use images from the registry directly with FauxRPC.

```shell
$ fauxrpc run --schema=buf.build/bufbuild/registry
```

This will start a fake version of the BSR API by downloading descriptors for [bufbuild/registry](https://buf.build/bufbuild/registry) from the BSR and using them with FauxRPC. Very meta.

## Buf Modules
If you give FauxRPC a path to a [Buf Module](https://buf.build/docs/cli/modules-workspaces/), it will automatically generate descriptors and use those, seamlessly. Note that you must also have [the buf CLI](https://buf.build/product/cli) installed and available in the path.

## Multiple Sources
You can define `--schema` as many times as needed. OpenAPI and Protobuf sources can be mixed in one server:
```shell
$ fauxrpc run --schema=./openapi.yaml --schema=./example.binpb
```

All the same inputs described above also work for [fauxrpc curl](/docs/client/fauxrpc-curl/).
