---
title: 'Protovalidate'
weight: 60
slug: protovalidate
---

FauxRPC uses [protovalidate](https://github.com/bufbuild/protovalidate) annotations to generate better fake data. The best way to understand is by showing.

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
