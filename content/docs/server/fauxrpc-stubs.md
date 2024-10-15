---
title: 'fauxrpc stub'
weight: 50
slug: fauxrpc-stub
---

FauxRPC's stubs allow you to specify what a response should look like. This is super useful for testing and setting up specific situations to happen during tests.

## Defining Stubs
Now, FauxRPC's default behavior is fun. It will generate [fake data](/docs/faking-data/) that probably looks a little silly. The data looks less silly if you [leverage protovalidate](/docs/faking-data/protovalidate/) to set constraints and examples. FauxRPC will automatically use those annotation to give more realistic data. This can be pretty decent, but there are many cases where this isn't good enough. You often times want to specify an exact response. Now you can do this with the `fauxrpc stubs` commands.

### Adding a stub
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
