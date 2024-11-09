---
title: 'Protovalidate'
slug: protovalidate
---

FauxRPC uses [protovalidate](https://github.com/bufbuild/protovalidate) annotations to generate better fake data. The best way to understand is by showing.

## Making better fake data

Now let's see what this looks like with protovalidate constraints:
```protobuf
syntax = "proto3";

import "buf/validate/validate.proto";

package greet.v1;

message GreetRequest {
  string name = 1 [(buf.validate.field).string = {min_len: 3, max_len: 100}];
}

message GreetResponse {
  string greeting = 1 [(buf.validate.field).string.example = "Hello, user!"];
}

service GreetService {
  rpc Greet(GreetRequest) returns (GreetResponse) {}
}
```

With this new protobuf file, this is what FauxRPC might output now:

```json
{
  "greeting": "Hello, user!"
}
```
This shows how protovalidate constraints enable FauxRPC to generate more realistic and contextually relevant fake data, aligning it closer to the expected behavior of your actual services.

Below is a `buf.registry.owner.v1.User` message defined in protobuf. This example comes from [the protobuf registry service](https://buf.build/bufbuild/registry/docs/main:buf.registry.owner.v1#buf.registry.owner.v1.User) for buf.build. Notice that the id field has protovalidate annotation saying that it's actually expected to have a `tuuid` format. TUUIDs are UUIDs but without the dashes. So a UUID of `7598a4f5-9a6a-4c68-9683-de54e21fb466` would look like `7598a4f59a6a4c689683de54e21fb466` in the TUUID format. Also, there's a name field with a minimum/maximum length and a regex pattern.

```protobuf
// A name uniquely identifies a User, however name is mutable.
message User {
  option (buf.registry.priv.extension.v1beta1.message).response_only = true;
  // The id for the User.
  string id = 1 [
    (buf.validate.field).required = true,
    (buf.validate.field).string.tuuid = true
  ];
  // The time the User was created.
  google.protobuf.Timestamp create_time = 2 [(buf.validate.field).required = true];
  // The last time the User was updated.
  google.protobuf.Timestamp update_time = 3 [(buf.validate.field).required = true];
  // The name of the User.
  //
  // A name uniquely identifies a User, however name is mutable.
  string name = 4 [
    (buf.validate.field).required = true,
    (buf.validate.field).string.min_len = 2,
    (buf.validate.field).string.max_len = 32,
    (buf.validate.field).string.pattern = "^[a-z][a-z0-9-]*[a-z0-9]$"
  ];
  // The type of the User.
  UserType type = 5 [
    (buf.validate.field).required = true,
    (buf.validate.field).enum.defined_only = true
  ];
  // The state of the User.
  UserState state = 6 [
    (buf.validate.field).required = true,
    (buf.validate.field).enum.defined_only = true
  ];
  // The configurable description of the User.
  string description = 7 [(buf.validate.field).string.max_len = 350];
  // The configurable URL that represents the homepage for a User.
  string url = 8 [
    (buf.validate.field).string.uri = true,
    (buf.validate.field).string.max_len = 255,
    (buf.validate.field).ignore = IGNORE_IF_UNPOPULATED
  ];
  // The verification status of the User.
  UserVerificationStatus verification_status = 9 [
    (buf.validate.field).required = true,
    (buf.validate.field).enum.defined_only = true
  ];
}
```

Now let's see the output that FauxRPC will give for this type:

```shell
$ buf build buf.build/bufbuild/registry -o registry.binpb
$ fauxrpc generate --schema=./registry.binpb --target=buf.registry.owner.v1.ListUsersResponse | jq
{
  "nextPageToken": "Microdosing.",
  "users": [
    {
      "id": "bc68e48242b646ca987cac9458460076",
      "createTime": "1925-11-15T00:18:39.375868764Z",
      "updateTime": "1907-10-26T13:51:30.881098122Z",
      "name": "s0",
      "type": "USER_TYPE_BOT",
      "description": "Freegan.",
      "url": "https://www.internationalreinvent.biz/enable"
    },
    {
      "id": "202ef6bec6c44aaab4ca3228045ed83c",
      "createTime": "2018-05-14T11:13:41.291735818Z",
      "updateTime": "1996-10-05T01:18:29.103348540Z",
      "name": "ske65l",
      "description": "Tilde.",
      "url": "https://www.producte-tailers.info/maximize",
      "verificationStatus": "USER_VERIFICATION_STATUS_VERIFIED"
    },
    {
      "id": "c294c494487e40e8a84ad0144e28b29c",
      "createTime": "2022-04-23T11:41:06.875834735Z",
      "updateTime": "1986-05-11T01:39:54.254076716Z",
      "name": "hliho",
      "type": "USER_TYPE_BOT",
      "state": "USER_STATE_INACTIVE",
      "description": "Master skateboard.",
      "url": "https://www.centralmorph.biz/transform/action-items/24-365",
      "verificationStatus": "USER_VERIFICATION_STATUS_VERIFIED"
    }
  ]
}
```
Hopefully, this gives you a good idea of what the output might look like. The better your validation rules, the better the FauxRPC data will be.

## Validation Rules
When running the FauxRPC server, these protovalidate constraints will actually be enforced (starting in [v0.1.2](https://github.com/sudorandom/fauxrpc/releases/tag/v0.1.2)). Let's see what that looks like!

Taking the `greet.v1.GreetService` from above, there's a little bit more work we have to do because running FauxRPC. This is because we used protovalidate, which is a dependency. Usually dependencies are hard to deal with in protobuf because the default protobuf tooling just doesn't handle it. However, the Buf CLI can help us here. Create a new `buf.yaml` file in the same directory as `greet.proto` with the contents:

```yaml
version: v2
deps:
  - buf.build/bufbuild/protovalidate
```

Now we're just a few commands away from having a mock service running:

```shell
# Get Buf CLI to pull our new dependency
$ buf dep update
# Build our protobuf (and dependencies) into a protobuf "image"
$ buf build . -o greet.binpb
# Start FauxRPC with this image
$ fauxrpc run --schema=greet.binpb
```

Now let's try some requests against this service that we just created:
```shell
$ buf curl --http2-prior-knowledge -d '{}' http://127.0.0.1:6660/greet.v1.GreetService/Greet
{
   "code": "invalid_argument",
   "message": "validation error:\n - name: value length must be at least 3 characters [string.min_len]",
   "details": [
      {
         "type": "buf.validate.Violations",
         "value": "CkIKBG5hbWUSDnN0cmluZy5taW5fbGVuGip2YWx1ZSBsZW5ndGggbXVzdCBiZSBhdCBsZWFzdCAzIGNoYXJhY3RlcnM",
         "debug": {
            "violations": [
               {
                  "fieldPath": "name",
                  "constraintId": "string.min_len",
                  "message": "value length must be at least 3 characters"
               }
            ]
         }
      }
   ]
}
```

Wow! This error tells us the exact issue with our request. The `name` field needs to be at least 3 characters now. Our empty object `{}` isn't good enough. We need a name and it needs to have between 3 and 20 characters. Let's try once more:

```shell
$ buf curl --http2-prior-knowledge -d '{"name": "Bob"}' http://127.0.0.1:6660/greet.v1.GreetService/Greet
{
  "greeting": "Hello, user!"
}
```

So while you are developing a frontend or some other kind of integration using a FauxRPC implementation, you will get real, high quality input validation feedback, so long as you use protovalidate.
