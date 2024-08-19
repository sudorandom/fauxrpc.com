---
title: 'Getting Started'
weight: 10
slug: getting-started
---


## How it Works
FauxRPC leverages your Protobuf definitions to generate fake services that mimic the behavior of real ones. You can easily configure the fake data returned, allowing you to simulate various scenarios and edge cases. It takes in `*.proto` files or protobuf descriptors (in binpb, json, txtpb, yaml formats), then it automatically starts up a server that can speak gRPC/gRPC-Web/Connect and REST (as long as there are `google.api.http` annotations defined). Descriptors contain all of the information found in a set of `.proto` files. You can generate them with `protoc` or the `buf build` command.

![](</diagram.svg>)

## Status: Alpha
This project is just starting out. I plan to add a lot of things that make this tool actually usable in more situations.

- Service for adding/updating/removing stub responses with a CLI to add/remove/replace these stubs
- Configuration file
- BSR Support (this is a 'maybe' because using `buf build` to emit descriptors works well enough IMO)
- Better streaming support. FauxRPC does work with streaming calls but it only returns a single response
