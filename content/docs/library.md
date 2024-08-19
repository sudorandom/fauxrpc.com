---
title: 'Library'
weight: 60
slug: library
---

You can use FauxRPC as a Go library as well. [Click here](https://pkg.go.dev/github.com/sudorandom/fauxrpc) for the most up-to-date documentation.

The two most important functions are:
- **[fauxrpc.SetDataOnMessage](https://pkg.go.dev/github.com/sudorandom/fauxrpc#SetDataOnMessage)**: When you have a concrete protobuf message type and want it populated with fake data.
- **[fauxrpc.NewMessage](https://pkg.go.dev/github.com/sudorandom/fauxrpc#NewMessage)**: When you have a MessageDescriptor and just want a message. This is typical for more dynamic environments where you may not have generated types to use.

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
    msg := &elizav1.SayResponse{}
		fauxrpc.SetDataOnMessage(msg)
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
		msg := &ownerv1.Owner{}
		fauxrpc.SetDataOnMessage(msg)
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

### dynamicpb
Sometimes you don't want to have the compiled protobuf code because you want the proxy or service to be more dynamic. In that case, you typically would use [dynamicpb](https://pkg.go.dev/google.golang.org/protobuf/types/dynamicpb). FauxRPC works there, too. Here's an example:

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

There are other functions to help generate specific fields. Please reference the [reference documentation](https://pkg.go.dev/github.com/sudorandom/fauxrpc) for more information on that.