package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"syscall/js"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/bufbuild/protocompile"
	"github.com/bufbuild/protocompile/reporter"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/sudorandom/fauxrpc"
	"github.com/sudorandom/fauxrpc/private/openapi/generator"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"

	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	_ "google.golang.org/protobuf/types/known/anypb"
	_ "google.golang.org/protobuf/types/known/durationpb"
	_ "google.golang.org/protobuf/types/known/emptypb"
	_ "google.golang.org/protobuf/types/known/structpb"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	_ "google.golang.org/protobuf/types/known/wrapperspb"
)

func main() {
	fmt.Println("WASM Proto Parser Initialized")
	js.Global().Set("parseProto", js.FuncOf(parseProto))
	js.Global().Set("generateFakeData", js.FuncOf(generateFakeData))
	js.Global().Set("formatPrototext", js.FuncOf(formatPrototext))
	js.Global().Set("listMessages", js.FuncOf(listMessages))
	js.Global().Set("parseOpenAPI", js.FuncOf(parseOpenAPI))
	js.Global().Set("generateOpenAPIFakeData", js.FuncOf(generateOpenAPIFakeData))
	js.Global().Set("listOpenAPIOperations", js.FuncOf(listOpenAPIOperations))
	select {}
}

func listMessages(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return js.ValueOf(map[string]any{"error": "missing arguments: fileDescriptorSet"})
	}

	fdsBytes := make([]byte, args[0].Get("length").Int())
	js.CopyBytesToGo(fdsBytes, args[0])

	fds := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(fdsBytes, fds); err != nil {
		return js.ValueOf(map[string]any{"error": "failed to unmarshal file descriptor set: " + err.Error()})
	}

	files, err := protodesc.NewFiles(fds)
	if err != nil {
		return js.ValueOf(map[string]any{"error": "failed to create protoregistry: " + err.Error()})
	}

	var names []any
	files.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		path := fd.Path()
		if strings.HasPrefix(path, "google/") || strings.HasPrefix(path, "buf/validate/") {
			return true
		}
		msgs := fd.Messages()
		for i := 0; i < msgs.Len(); i++ {
			names = append(names, string(msgs.Get(i).FullName()))
		}
		return true
	})

	return js.ValueOf(names)
}

func formatPrototext(this js.Value, args []js.Value) any {
	if len(args) < 3 {
		return js.ValueOf(map[string]any{"error": "missing arguments: messageName, fileDescriptorSet, jsonString"})
	}

	messageName := args[0].String()
	fdsBytes := make([]byte, args[1].Get("length").Int())
	js.CopyBytesToGo(fdsBytes, args[1])
	jsonString := args[2].String()

	fds := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(fdsBytes, fds); err != nil {
		return js.ValueOf(map[string]any{"error": "failed to unmarshal file descriptor set: " + err.Error()})
	}

	files, err := protodesc.NewFiles(fds)
	if err != nil {
		return js.ValueOf(map[string]any{"error": "failed to create protoregistry: " + err.Error()})
	}

	desc, err := files.FindDescriptorByName(protoreflect.FullName(messageName))
	if err != nil {
		return js.ValueOf(map[string]any{"error": "failed to find message descriptor: " + err.Error()})
	}

	messageDesc, ok := desc.(protoreflect.MessageDescriptor)
	if !ok {
		return js.ValueOf(map[string]any{"error": "descriptor is not a message: " + messageName})
	}

	msg := dynamicpb.NewMessage(messageDesc)
	
	// Unmarshal options to ignore unknown fields just in case
	unmarshalOpts := protojson.UnmarshalOptions{DiscardUnknown: true}
	err = unmarshalOpts.Unmarshal([]byte(jsonString), msg)
	if err != nil {
		return js.ValueOf(map[string]any{"error": "failed to unmarshal json: " + err.Error()})
	}

	textOpts := prototext.MarshalOptions{
		Multiline: true,
		Indent:    "  ",
	}
	textBytes, err := textOpts.Marshal(msg)
	if err != nil {
		return js.ValueOf(map[string]any{"error": "failed to marshal prototext: " + err.Error()})
	}

	return js.ValueOf(string(textBytes))
}

func generateFakeData(this js.Value, args []js.Value) any {
	if len(args) < 2 {
		return js.ValueOf(map[string]any{"error": "missing arguments: messageName, fileDescriptorSet"})
	}

	messageName := args[0].String()
	fdsBytes := make([]byte, args[1].Get("length").Int())
	js.CopyBytesToGo(fdsBytes, args[1])

	fds := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(fdsBytes, fds); err != nil {
		return js.ValueOf(map[string]any{"error": "failed to unmarshal file descriptor set: " + err.Error()})
	}

	files, err := protodesc.NewFiles(fds)
	if err != nil {
		return js.ValueOf(map[string]any{"error": "failed to create protoregistry: " + err.Error()})
	}

	desc, err := files.FindDescriptorByName(protoreflect.FullName(messageName))
	if err != nil {
		return js.ValueOf(map[string]any{"error": "failed to find message descriptor: " + err.Error()})
	}

	messageDesc, ok := desc.(protoreflect.MessageDescriptor)
	if !ok {
		return js.ValueOf(map[string]any{"error": "descriptor is not a message: " + messageName})
	}

	msg, err := fauxrpc.NewMessage(messageDesc, fauxrpc.GenOptions{
		Faker: gofakeit.New(uint64(time.Now().UnixNano())),
	})
	if err != nil {
		return js.ValueOf(map[string]any{"error": "failed to generate fake data: " + err.Error()})
	}

	jsonBytes, err := protojson.MarshalOptions{EmitUnpopulated: true, Indent: "  "}.Marshal(msg)
	if err != nil {
		return js.ValueOf(map[string]any{"error": "failed to marshal fake data to JSON: " + err.Error()})
	}

	return js.ValueOf(string(jsonBytes))
}

type compilationError struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Col     int    `json:"col"`
	Offset  int    `json:"offset"`
	Message string `json:"message"`
}

func parseProto(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return js.ValueOf(map[string]any{"error": "missing arguments"})
	}

	var files map[string]string
	err := json.Unmarshal([]byte(args[0].String()), &files)
	if err != nil {
		return js.ValueOf(map[string]any{"error": "invalid json: " + err.Error()})
	}

	var compilationErrors []compilationError
	rep := reporter.NewReporter(func(err reporter.ErrorWithPos) error {
		pos := err.GetPosition()
		compilationErrors = append(compilationErrors, compilationError{
			File:    pos.Filename,
			Line:    pos.Line,
			Col:     pos.Col,
			Offset:  pos.Offset,
			Message: err.Unwrap().Error(),
		})
		return nil // Continue to find more errors
	}, nil)

	compiler := protocompile.Compiler{
		Resolver: protocompile.CompositeResolver{
			&protocompile.SourceResolver{
				Accessor: func(path string) (io.ReadCloser, error) {
					content, ok := files[path]
					if !ok {
						return nil, fmt.Errorf("file not found: %s", path)
					}
					return io.NopCloser(strings.NewReader(content)), nil
				},
			},
			protocompile.ResolverFunc(func(path string) (protocompile.SearchResult, error) {
				fd, err := protoregistry.GlobalFiles.FindFileByPath(path)
				if err != nil {
					return protocompile.SearchResult{}, err
				}
				return protocompile.SearchResult{
					Proto: protodesc.ToFileDescriptorProto(fd),
				}, nil
			}),
		},
		Reporter: rep,
	}

	// Compile all files provided in the input
	var fileNames []string
	for name := range files {
		fileNames = append(fileNames, name)
	}

	ctx := context.Background()
	fds, err := compiler.Compile(ctx, fileNames...)
	if err != nil {
		if len(compilationErrors) > 0 {
			errorsJson, _ := json.Marshal(compilationErrors)
			return js.ValueOf(map[string]any{"compilationErrors": string(errorsJson)})
		}

		var errWithPos reporter.ErrorWithPos
		if errors.As(err, &errWithPos) {
			pos := errWithPos.GetPosition()
			e := compilationError{
				File:    pos.Filename,
				Line:    pos.Line,
				Col:     pos.Col,
				Offset:  pos.Offset,
				Message: errWithPos.Unwrap().Error(),
			}
			errorsJson, _ := json.Marshal([]compilationError{e})
			return js.ValueOf(map[string]any{"compilationErrors": string(errorsJson)})
		}
		return js.ValueOf(map[string]any{"error": err.Error()})
	}

	fdpSet := &descriptorpb.FileDescriptorSet{}
	seen := make(map[string]bool)
	var addFile func(fd protoreflect.FileDescriptor)
	addFile = func(fd protoreflect.FileDescriptor) {
		if seen[fd.Path()] {
			return
		}
		seen[fd.Path()] = true
		for i := 0; i < fd.Imports().Len(); i++ {
			addFile(fd.Imports().Get(i).FileDescriptor)
		}
		fdpSet.File = append(fdpSet.File, protodesc.ToFileDescriptorProto(fd))
	}

	for _, fd := range fds {
		addFile(fd)
	}

	bytes, err := proto.Marshal(fdpSet)
	if err != nil {
		return js.ValueOf(map[string]any{"error": "marshal error: " + err.Error()})
	}

	// Return as Uint8Array
	uint8Array := js.Global().Get("Uint8Array").New(len(bytes))
	js.CopyBytesToJS(uint8Array, bytes)

	return js.ValueOf(map[string]any{"fileDescriptorSet": uint8Array})
}

var openapiDocCache *openapi3.T

func parseOpenAPI(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return js.ValueOf(map[string]any{"error": "missing arguments: spec content"})
	}

	content := args[0].String()
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	doc, err := loader.LoadFromData([]byte(content))
	if err != nil {
		return js.ValueOf(map[string]any{"error": "failed to parse openapi spec: " + err.Error()})
	}

	ctx := context.Background()
	if err := doc.Validate(ctx); err != nil {
		return js.ValueOf(map[string]any{"error": "invalid openapi spec: " + err.Error()})
	}

	openapiDocCache = doc

	var ops []any
	if doc.Paths != nil {
		for path, item := range doc.Paths.Map() {
			if item == nil {
				continue
			}
			for method, op := range item.Operations() {
				if op == nil {
					continue
				}
				opID := op.OperationID
				if opID == "" {
					opID = fmt.Sprintf("%s %s", method, path)
				} else {
					opID = fmt.Sprintf("%s (%s %s)", opID, method, path)
				}
				ops = append(ops, opID)
			}
		}
	}
	sort.Slice(ops, func(i, j int) bool {
		return fmt.Sprint(ops[i]) < fmt.Sprint(ops[j])
	})

	return js.ValueOf(map[string]any{"operations": ops})
}

func listOpenAPIOperations(this js.Value, args []js.Value) any {
	if openapiDocCache == nil || openapiDocCache.Paths == nil {
		return js.ValueOf([]any{})
	}

	var ops []any
	for path, item := range openapiDocCache.Paths.Map() {
		if item == nil {
			continue
		}
		for method, op := range item.Operations() {
			if op == nil {
				continue
			}
			opID := op.OperationID
			if opID == "" {
				opID = fmt.Sprintf("%s %s", method, path)
			} else {
				opID = fmt.Sprintf("%s (%s %s)", opID, method, path)
			}
			ops = append(ops, opID)
		}
	}
	sort.Slice(ops, func(i, j int) bool {
		return fmt.Sprint(ops[i]) < fmt.Sprint(ops[j])
	})

	return js.ValueOf(ops)
}

func generateOpenAPIFakeData(this js.Value, args []js.Value) any {
	if openapiDocCache == nil {
		return js.ValueOf(map[string]any{"error": "no valid openapi spec loaded"})
	}

	if len(args) < 1 {
		return js.ValueOf(map[string]any{"error": "missing argument: operationTarget"})
	}

	target := args[0].String()
	walker := generator.NewWalker(false)

	for path, item := range openapiDocCache.Paths.Map() {
		if item == nil {
			continue
		}
		for method, op := range item.Operations() {
			if op == nil {
				continue
			}
			opID := op.OperationID
			opLabel := opID
			if opLabel == "" {
				opLabel = fmt.Sprintf("%s %s", method, path)
			} else {
				opLabel = fmt.Sprintf("%s (%s %s)", opID, method, path)
			}

			if target == opLabel || target == opID || target == fmt.Sprintf("%s %s", method, path) {
				statusCode, headers, payload, err := walker.GenerateFromOperation(method, path, opID, op, 5)
				if err != nil {
					return js.ValueOf(map[string]any{"error": "failed to generate payload: " + err.Error()})
				}

				jsonBytes, err := json.MarshalIndent(payload, "", "  ")
				if err != nil {
					return js.ValueOf(map[string]any{"error": "failed to encode JSON: " + err.Error()})
				}

				resMap := map[string]any{
					"status":  statusCode,
					"payload": string(jsonBytes),
				}
				if len(headers) > 0 {
					resMap["headers"] = headers
				}
				return js.ValueOf(resMap)
			}
		}
	}

	return js.ValueOf(map[string]any{"error": "operation not found: " + target})
}
