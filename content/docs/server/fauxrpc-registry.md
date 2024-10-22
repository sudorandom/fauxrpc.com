---
title: 'fauxrpc registry'
weight: 50
slug: fauxrpc-registry
---

The "registry" refers to the schema that FauxRPC knows about. Since this a protobuf/gRPC focused tool "schema" comes in the from of [protobuf descriptors](https://buf.build/docs/reference/descriptors). Descriptors can come from many sources but the ways to input them are detailed further in the [inputs page](/docs/server/inputs/).

### Adding new schema
Schema can be added using the `--schema` option whenever starting the server but you can actually change the schema in runtime as well using the `fauxrpc registry add` command.

Let's launch a FauxRPC server and populate it with schema to see this working in action.
```shell
$ fauxrpc run --empty
FauxRPC (dev (none) @ unknown; go1.23.2) - 0 services loaded
Listening on http://127.0.0.1:6660
OpenAPI documentation: http://127.0.0.1:6660/fauxrpc/openapi.html

Example Commands:
$ buf curl --http2-prior-knowledge http://127.0.0.1:6660 --list-methods
$ buf curl --http2-prior-knowledge http://127.0.0.1:6660/[METHOD_NAME]
Server started.
```
Now that we have a server, we can work on another terminal to populate it with schema.

```shell
# output the "image" in the form of protobuf encoded file descriptors:
$ buf build buf.build/connectrpc/eliza -o eliza.binpb
# give it to the FauxRPC service
$ fauxrpc registry add
```
Okay, that didn't give any output. Now what? Well, now FauxRPC knows about the "Eliza" service. We can now use existing tools that leverage these services. Since FauxRPC supports gRPC server reflection, we can also use that to discover which methods are available:

```shell
$ buf curl --http2-prior-knowledge http://127.0.0.1:6660 --list-methods
connectrpc.eliza.v1.ElizaService/Converse
connectrpc.eliza.v1.ElizaService/Introduce
connectrpc.eliza.v1.ElizaService/Say
```

Sweet, we have three methods available to use!
```shell
buf curl --http2-prior-knowledge -d '{"sentence": "Hello World!"}' http://127.0.0.1:6660/connectrpc.eliza.v1.ElizaService/Say
{
  "sentence": "Lumbersexual."
}
```
And it works! It's outputting absolute nonsense! There's 

### Removing schema

```shell
$ fauxrpc registry remove-all
```
POOF! Now our FauxRPC is a blank slate again. We can verify this by listing the available methods again:
```shell
buf curl --http2-prior-knowledge http://127.0.0.1:6660 --list-methods
```
