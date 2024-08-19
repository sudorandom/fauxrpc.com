---
title: 'Introduction'
weight: 10
slug: intro
aliases:
  - "/docs/"
---

## FauxRPC
**FauxRPC** is a powerful tool that empowers you to accelerate development and testing by effortlessly generating fake implementations of gRPC, gRPC-Web, Connect, and REST services. If you have a protobuf-based workflow, this tool could help.

## How it Works
FauxRPC leverages your Protobuf definitions to generate fake services that mimic the behavior of real ones. You can easily configure the fake data returned, allowing you to simulate various scenarios and edge cases. It takes in `*.proto` files or protobuf descriptors (in binpb, json, txtpb, yaml formats), then it automatically starts up a server that can speak gRPC/gRPC-Web/Connect and REST (as long as there are `google.api.http` annotations defined). Descriptors contain all of the information found in a set of `.proto` files. You can generate them with `protoc` or the `buf build` command.

![](</diagram.svg>)

## Status: Alpha
This project is just starting out. I plan to add a lot of things that make this tool actually usable in more situations.

Keep track of the progress by following [issues](https://github.com/sudorandom/fauxrpc/issues) and [pull requests](https://github.com/sudorandom/fauxrpc/pulls) on GitHub.
