---
title: 'Faking Data'
weight: 50
slug: faking-data
---

## So what does the fake data look like?
You might be wondering what actual responses look like. FauxRPC's fake data generation is continually improving so these details might change as time goes on. It uses a library called [fakeit](https://github.com/brianvoe/gofakeit) to generate fake data. Because protobufs have pretty well-defined types, we can easily generate data that technically matches the types. This works well for most use cases, but FauxRPC tries to be a little bit better. If you annotate your protobuf files with [protovalidate](https://github.com/bufbuild/protovalidate) constraints, FauxRPC will try its best to generate data that matches these constraints. Let's look at some examples!

### With no annotations
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

With FauxRPC, you will get any kind of word, so it might look like this:
```json
{
  "greeting": "Poutine."
}
```
This is fine, but for this RPC, we know a bit more about the type being returned. We know that it sends a greeting back that looks like "Hello, [name]". So here's what the same protobuf file might look like with protovalidate constraints:

### With protovalidate
You can create more realistic fake data just by implementing protovalidate rules. See more on this topic [here](/docs/protovalidate/).
