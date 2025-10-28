---
title: 'fauxrpc stub'
weight: 50
slug: fauxrpc-stub
description: "Learn how to use the `fauxrpc stub` command to define specific, static responses for your RPC calls."
icon: "extension"
---

FauxRPC's stubs allow you to specify what a response should look like. This is super useful for testing and setting up specific situations to happen during tests.

```shell
Usage: fauxrpc stub <command> [flags]

Contains stub commands

Flags:
  -h, --help                Show context-sensitive help.
  -l, --log-level="info"    Set the logging level (debug|info|warn|error)
      --version             Print version information and quit

Commands:
  stub add           Adds a new stub response by method or type
  stub list          List all registered stubs
  stub get           Get a registered stub
  stub remove        Remove a registered stub
  stub remove-all    Remove all stubs
```

## Defining Stubs
Now, FauxRPC's default behavior is fun. It will generate [fake data](/docs/faking-data/) that probably looks a little silly. The data looks less silly if you [leverage protovalidate](/docs/protovalidate/) to set constraints and examples. FauxRPC will automatically use those annotation to give more realistic data. This can be pretty decent, but there are many cases where this isn't good enough. You often times want to specify an exact response. Now you can do this with the `fauxrpc stubs` commands.

### Adding a stub
```shell
Usage: fauxrpc stub add <target> [flags]

Adds a new stub response by method or type

Arguments:
  <target>    Protobuf method or type

Flags:
  -h, --help                            Show context-sensitive help.
  -l, --log-level="info"                Set the logging level (debug|info|warn|error)
      --version                         Print version information and quit

  -a, --addr="http://127.0.0.1:6660"    Address to bind to.
      --id=STRING                       ID to give this particular mock response, will be a random string if one isn't given
      --json=STRING                     Protobuf method or type
      --error-message=STRING            Message to return with the error
      --error-code=ERROR-CODE           gRPC Error code to return
      --active-if=STRING                CEL expression that must be true before this mock is used.
      --priority=0                      Priority from 0-100 (higher is more preferred)
      --cel=STRING                      CEL expression
```

```shell
# We are starting fauxrpc with our schema, but this schema could also be added dynamically.
$ fauxrpc run --schema=eliza.binpb
# Set a stub for the connectrpc.eliza.v1.ElizaService/Say method 
$ fauxrpc stub add connectrpc.eliza.v1.ElizaService/Say --json '{"sentence": "Test from the CLI"}'
```
This `fauxrpc stub add` command defined what the FauxRPC server should return when requesting `connectrpc.eliza.v1.ElizaService/Say`. Let's try it!
```shell
$ buf curl --http2-prior-knowledge -d '{}' http://127.0.0.1:6660/connectrpc.eliza.v1.ElizaService/Say
{
  "sentence": "Test from the CLI"
}
```
Yep! the service is responding correctly!

### Removing Stubs
Similar to the registry, there's also a way to clear out all of the registered stubs. Here's what that looks like:
```shell
$ fauxrpc stub remove-all

# Now the method is back to the nonsense you know and love from FauxRPC!
$ buf curl --http2-prior-knowledge -d '{}' http://127.0.0.1:6660/connectrpc.eliza.v1.ElizaService/Say
{
  "sentence": "Small batch pickled brunch."
}
```

Also be sure to check out the `list`, `get` and `remove` commands. These can be used to manage specific stubbed responses so you don't have to rebuild the world again after calling remove-all.
