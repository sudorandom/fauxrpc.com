module github.com/sudorandom/fauxrpc/wasm-fauxrpc

go 1.26.3

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.11-20260709200747-435963d16310.1
	github.com/brianvoe/gofakeit/v7 v7.15.0
	github.com/bufbuild/protocompile v0.14.1
	github.com/getkin/kin-openapi v0.143.0
	github.com/sudorandom/fauxrpc v0.20.1
	google.golang.org/genproto/googleapis/api v0.0.0-20260706201446-f0a921348800
	google.golang.org/protobuf v1.36.11
)

replace github.com/sudorandom/fauxrpc => ../..

require (
	buf.build/gen/go/bufbuild/registry/connectrpc/go v1.20.0-20260713175918-10d915f5b43b.1 // indirect
	buf.build/gen/go/bufbuild/registry/protocolbuffers/go v1.36.11-20260713175918-10d915f5b43b.1 // indirect
	buf.build/gen/go/grpc/grpc/connectrpc/go v1.20.0-20260331211127-1730f7242d0f.1 // indirect
	buf.build/gen/go/grpc/grpc/protocolbuffers/go v1.36.11-20260331211127-1730f7242d0f.1 // indirect
	buf.build/go/protovalidate v1.2.0 // indirect
	buf.build/go/protoyaml v0.7.0 // indirect
	cel.dev/expr v0.25.2 // indirect
	connectrpc.com/connect v1.20.0 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/go-openapi/jsonpointer v0.22.5 // indirect
	github.com/go-openapi/swag/jsonname v0.25.5 // indirect
	github.com/google/cel-go v0.29.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/oasdiff/yaml v0.1.1 // indirect
	github.com/oasdiff/yaml3 v0.0.14 // indirect
	github.com/santhosh-tekuri/jsonschema/v6 v6.0.2 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/exp v0.0.0-20260709172345-9ea1abe57597 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260727163830-6c54dddc4772 // indirect
)
