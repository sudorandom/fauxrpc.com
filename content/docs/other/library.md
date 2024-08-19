---
title: 'Library'
weight: 60
slug: library
---

You can use FauxRPC as a Go library as well. [Click here](https://pkg.go.dev/github.com/sudorandom/fauxrpc) for the most up-to-date documentation.

Here's how you can generate a message filled with fake data, using the Eliza service example.
```go
package main

import (
	"fmt"
	"log"

	elizav1 "buf.build/gen/go/connectrpc/eliza/protocolbuffers/go/connectrpc/eliza/v1"
	"github.com/sudorandom/fauxrpc"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
    msg := fauxrpc.NewMessage(elizav1.File_connectrpc_eliza_v1_eliza_proto.Messages().ByName("SayResponse"))
    b, err := protojson.MarshalOptions{Indent: "  "}.Marshal(msg)
    if err != nil {
        log.Fatalf("err: %s", err)
    }
    fmt.Println(string(b))
}
```

Example output:
```json
{
  "sentence": "Jean shorts."
}
```

This text will be randomly generated each time. That's nice, but Here's one that's a bit more complex.
```go
package main

import (
	"fmt"
	"log"

	ownerv1 "buf.build/gen/go/bufbuild/registry/protocolbuffers/go/buf/registry/owner/v1"
	"github.com/sudorandom/fauxrpc"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
    msg := fauxrpc.NewMessage(ownerv1.File_buf_registry_owner_v1_owner_proto.Messages().ByName("Owner"))
    b, err := protojson.MarshalOptions{Indent: "  "}.Marshal(msg)
    if err != nil {
        log.Fatalf("err: %s", err)
    }
    fmt.Println(string(b))
}
```

Example output:
```json
{
  "organization": {
    "id": "a4cf6166453b49d9811cfdc169c36354",
    "createTime": "1912-06-01T16:45:13.510830647Z",
    "updateTime": "2009-04-03T13:11:31.216229245Z",
    "name": "xm0r3",
    "description": "Godard selvage.",
    "url": "https://www.dynamicreintermediate.name/robust/seize/metrics/b2c",
    "verificationStatus": "ORGANIZATION_VERIFICATION_STATUS_OFFICIAL"
  }
}
```

There are functions to help generate specific fields. Please reference the [reference documentation](https://pkg.go.dev/github.com/sudorandom/fauxrpc) for more information on that.