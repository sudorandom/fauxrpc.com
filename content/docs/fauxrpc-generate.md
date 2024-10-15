---
title: 'fauxrpc generate'
weight: 50
slug: fauxrpc-generate
---

The `fauxrpc generate` command is a handy tool within FauxRPC that allows you to generate fake data for your Protobuf message types. This is useful for testing, seeding databases, or creating realistic (and often-times silly) examples for documentation. Here's a breakdown of the command and its options:

## Usage
```bash
Usage: fauxrpc generate --schema=SCHEMA,... --target=STRING [flags]

Generate fake data

Flags:
  -h, --help                 Show context-sensitive help.
  -l, --log-level="info"     Set the logging level (debug|info|warn|error)
      --version              Print version information and quit

      --schema=SCHEMA,...    The modules to use for the RPC schema. It can be protobuf descriptors (binpb, json, yaml),
                             a URL for reflection or a directory of descriptors.
      --target=STRING        Protobuf type
      --format="json"        Format to output
      --seed=SEED            Seed for random number generator
```

## Flags

- `-h, --help`: Displays help information for the command.
- `-l, --log-level`: Sets the logging level (`debug`, `info`, `warn`, `error`). Default is `info`.
- `--version`: Prints the FauxRPC version.
- `--schema`: Specifies the source for your Protobuf definitions. It accepts:
    - **Protobuf descriptors:** Paths to files containing Protobuf definitions (`.binpb`, `.json`, `.yaml`). This is the most common way to define your services.
    - **Reflection URL:** A URL of a live gRPC server that supports reflection. FauxRPC will query this server to understand its services and generate fake responses accordingly.
    - **Directory:** A path to a directory containing Protobuf descriptor files.
    * You can use this flag multiple times to combine services from different sources.
    * See the [Inputs](/docs/server/inputs) page for more on this.
- `--target`: **Required.** Defines the fully qualified Protobuf message type for which you want to generate fake data (e.g., `my.package.v1.MyMessage`).
- `--format`: Sets the output format for the generated data. Default is `json`.
- `--seed`: Provides a seed value for the random number generator. This ensures consistent output if you need to generate the same data multiple times.

**Example:**

```bash
# Build the Protobuf descriptors
buf build buf.build/connectrpc/eliza -o eliza.binpb 

# Generate fake data for the SayRequest message in JSON format
fauxrpc generate --schema=eliza.binpb --target=connectrpc.eliza.v1.SayRequest 
```

This example will produce a JSON output containing fake data for the `connectrpc.eliza.v1.SayRequest` message type, based on the definitions in the `eliza.binpb` file. Here's some examples of the output:

```shell
$ fauxrpc generate --schema=eliza.binpb --target=connectrpc.eliza.v1.SayRequest 
{"sentence":"Messenger bag."}

$ buf build buf.build/bufbuild/registry -o registry.binpb
$ fauxrpc generate --schema=registry.binpb --target=buf.registry.module.v1.Module | jq
{
  "id": "1abe706eea5a4826b0b8b21ed8caaaa0",
  "createTime": "2012-02-04T20:15:42.740555734Z",
  "updateTime": "1953-12-31T19:25:51.725429551Z",
  "name": "Plaid lumbersexual cold-pressed.",
  "ownerId": "0e64a0beee864e80b06e3ce28ef8ba84",
  "visibility": "MODULE_VISIBILITY_PUBLIC",
  "state": "MODULE_STATE_DEPRECATED",
  "description": "Paleo.",
  "url": "https://www.leadout-of-the-box.org/eyeballs/visionary/reinvent",
  "defaultLabelName": "Polaroid fixie."
}

$ fauxrpc generate --schema=registry.binpb --target=buf.registry.module.v1.Commit | jq
{
  "id": "ba448e58081e4eda9062a2e8ae2cfaae",
  "createTime": "2024-02-21T10:33:29.535690109Z",
  "ownerId": "a0d2c0b4bafa42cc9af8082a12548610",
  "moduleId": "525e560c6818442eaa868a2ea62a36fe",
  "digest": {
    "value": "WW91IHByb2JhYmx5IGhhdmVuJ3QgaGVhcmQgb2YgdGhlbSB0aWxkZSBiaWN5Y2xlIHJpZ2h0cy4="
  },
  "createdByUserId": "fe9e90c8706e4c7494a8ec8aaa847470",
  "sourceControlUrl": "https://www.directe-business.name/real-time/holistic/technologies"
}
```

**Key Takeaways:**

* `fauxrpc generate` helps you create realistic fake data for your Protobuf messages.
* You can control the output format and use a seed for reproducible results.
* It's a valuable tool for testing, prototyping, and various development scenarios where you need sample data.

Remember that you can [use protovalidate annotations](/docs/faking-data/protovalidate) to improve the quality of this data.
