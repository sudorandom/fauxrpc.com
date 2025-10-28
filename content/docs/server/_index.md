---
slug: "_index"
title: "Server"
weight: 40
description: "Learn how to use the FauxRPC server for testing and development."
icon: "dns"
---

FauxRPC is a handy tool that lets you quickly spin up a fake server that mimics a real gRPC server. This is super useful for testing and development when you don't have the actual backend service ready yet, or when you want to isolate your application from the real backend for focused testing.

## Inputs
FauxRPC allows you to easily mock gRPC servers using various inputs like Protobuf files, live gRPC servers, or Buf Schema Registry images. This makes it simple to test gRPC interactions without a real backend. See more on the [Inputs](/docs/server/inputs/) page.

## Multiprotocol Support
FauxRPC lets you easily create fake gRPC, gRPC-Web, Connect, and REST servers from different sources like Protobuf files, live gRPC servers, or Buf Schema Registry images. This allows you to test and develop your applications without needing a real backend. You can use tools like buf curl or grpcurl to interact with the fake server. See details on the [Multi-protocol Support](/docs/server/multi-protocol-support/) page.

## Dynamic Reconfiguration

FauxRPC allows you to send it protobuf descriptors and it will mimic the services contained within those descriptors. This allows you to set up fairly complex scenarios dynamically. This is all powered by an incredibly [simple gRPC service](https://github.com/sudorandom/fauxrpc/blob/main/proto/registry/v1/registry_service.proto) and easy to use with FauxRPC with the CLI.

See the documentation on [fauxrpc registry](/docs/server/fauxrpc-registry) and [fauxrpc stubs](/docs/server/fauxrpc-registry/) for more.