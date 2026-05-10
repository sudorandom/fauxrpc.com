module github.com/sudorandom/fauxrpc/wasm-fauxrpc

go 1.25.7

replace github.com/sudorandom/fauxrpc => ../../fauxrpc

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.11-20260415201107-50325440f8f2.1
	github.com/bufbuild/protocompile v0.14.1
	github.com/sudorandom/fauxrpc v0.0.0-00010101000000-000000000000
	google.golang.org/genproto/googleapis/api v0.0.0-20260504160031-60b97b32f348
	google.golang.org/protobuf v1.36.11
)

require (
	buf.build/gen/go/bufbuild/registry/connectrpc/go v1.19.1-20251027152159-f1066ce064ca.2 // indirect
	buf.build/gen/go/bufbuild/registry/protocolbuffers/go v1.36.10-20251027152159-f1066ce064ca.1 // indirect
	buf.build/gen/go/grpc/grpc/connectrpc/go v1.19.1-20260203201457-e126be52bace.2 // indirect
	buf.build/gen/go/grpc/grpc/protocolbuffers/go v1.36.10-20260203201457-e126be52bace.1 // indirect
	buf.build/go/protovalidate v1.1.2 // indirect
	buf.build/go/protoyaml v0.6.0 // indirect
	cel.dev/expr v0.25.1 // indirect
	connectrpc.com/connect v1.19.1 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/brianvoe/gofakeit/v7 v7.8.1 // indirect
	github.com/google/cel-go v0.27.0 // indirect
	golang.org/x/exp v0.0.0-20260218203240-3dfff04db8fa // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260427160629-7cedc36a6bc4 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
