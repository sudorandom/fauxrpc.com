
run: build-wasm
	hugo server --buildDrafts --buildFuture --minify -p 1414

build-wasm:
	cd wasm-fauxrpc && GOOS=js GOARCH=wasm go build -o ../static/fauxrpc.wasm main.go
