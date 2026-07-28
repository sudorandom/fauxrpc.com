---
slug: "_index"
title: "Server"
weight: 40
description: "Learn how to use the FauxRPC server for testing and development."
icon: "dns"
---

FauxRPC quickly spins up fake Protobuf and OpenAPI services. This is useful when the real backend is not ready or when a test needs an isolated, controllable implementation.

## Inputs
FauxRPC accepts OpenAPI YAML/JSON, Protobuf descriptor sets, live gRPC reflection endpoints, and Buf Schema Registry images through the same `--schema` option. See [Inputs](/docs/server/inputs/) for all supported sources.

## OpenAPI Support

Serve schema-aware HTTP operations, validate requests, generate bodies and response headers, load conditional stubs, and browse interactive API documentation. See [OpenAPI Support](/docs/server/openapi/) for a complete guide.

## Multiprotocol Support
FauxRPC lets you easily create fake gRPC, gRPC-Web, Connect, and REST servers from different sources like Protobuf files, live gRPC servers, or Buf Schema Registry images. This allows you to test and develop your applications without needing a real backend. You can use tools like buf curl or grpcurl to interact with the fake server. See details on the [Multi-protocol Support](/docs/server/multi-protocol-support/) page.

## Dynamic Reconfiguration

FauxRPC allows you to send it protobuf descriptors and it will mimic the services contained within those descriptors. This allows you to set up fairly complex scenarios dynamically. This is all powered by an incredibly [simple gRPC service](https://github.com/sudorandom/fauxrpc/blob/main/proto/registry/v1/registry_service.proto) and easy to use with FauxRPC with the CLI.

See the documentation on [Manage Registry](/docs/server/fauxrpc-registry/) and [Manage Stubs via CLI](/docs/server/fauxrpc-stub/) for more.

## Upstream Proxying & Recording

FauxRPC can act as an intercepting proxy between your client and an upstream server. It forwards requests, intercepts `UNIMPLEMENTED` errors to serve mocks or fakes, and can record all passing traffic to structured JSON stub files for offline replay.

See the documentation on [Proxy & Record](/docs/server/proxy-and-record/) for more.
