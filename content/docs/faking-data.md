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
Now let's see what this looks like with protovalidate constraints:
```protobuf
syntax = "proto3";

import "buf/validate/validate.proto";

package greet.v1;

message GreetRequest {
  string name = 1 [(buf.validate.field).string = {min_len: 3, max_len: 100}];
}

message GreetResponse {
  string greeting = 1 [(buf.validate.field).string.pattern = "^Hello, [a-zA-Z]+$"];
}

service GreetService {
  rpc Greet(GreetRequest) returns (GreetResponse) {}
}
```

With this new protobuf file, this is what FauxRPC might output now:

```json
{
  "greeting": "Hello, TWXxF"
}
```
This shows how protovalidate constraints enable FauxRPC to generate more realistic and contextually relevant fake data, aligning it closer to the expected behavior of your actual services. As another example, I will show one of Buf's services used to manage users, [buf.registry.owner.v1.UserService](https://buf.build/bufbuild/registry/docs/main:buf.registry.owner.v1#buf.registry.owner.v1.UserService). Here's what the `UserRef` message looks like:

```protobuf
message UserRef {
  option (buf.registry.priv.extension.v1beta1.message).request_only = true;
  oneof value {
    option (buf.validate.oneof).required = true;
    // The id of the User.
    string id = 1 [(buf.validate.field).string.tuuid = true];
    // The name of the User.
    string name = 2 [(buf.validate.field).string = {
      min_len: 2
      max_len: 32
      pattern: "^[a-z][a-z0-9-]*[a-z0-9]$"
    }];
  }
}
```

So let's make our descriptors for this service, start the FauxRPC server and make our example request:

```shell
$ buf build buf.build/bufbuild/registry -o bufbuild.registry.binpb
$ fauxrpc run --schema=./bufbuild.registry.binpb
$ buf curl --http2-prior-knowledge http://127.0.0.1:6660/buf.registry.owner.v1.UserService/ListUsers
{
  "nextPageToken": "Food truck.",
  "users": [
    {
      "id": "c4468393f926400d8880a264df9c284a",
      "createTime": "2012-03-06T12:15:03.239463070Z",
      "updateTime": "1990-10-29T13:12:31.224347086Z",
      "name": "jexox",
      "type": "USER_TYPE_STANDARD",
      "description": "Tattooed taxidermy.",
      "url": "http://www.productexploit.name/synergies/target"
    },
    {
      "id": "0e4ca24f4ff54761b109daab0da1bea2",
      "createTime": "1955-05-16T02:37:30.643378679Z",
      "updateTime": "1923-08-28T04:28:43.330711919Z",
      "name": "ya0",
      "type": "USER_TYPE_STANDARD",
      "state": "USER_STATE_INACTIVE",
      "description": "Helvetica.",
      "url": "https://www.centralengage.info/markets/scale/e-commerce/exploit",
      "verificationStatus": "USER_VERIFICATION_STATUS_UNVERIFIED"
    }
  ]
}
```
Hopefully, this gives you a good idea of what the output might look like. The better your validation rules, the better the FauxRPC data will be.