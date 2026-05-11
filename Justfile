
run: build-wasm
	hugo server --buildDrafts --buildFuture --minify -p 1414

build-wasm:
	@if [ -f static/fauxrpc.wasm ] && [ -z "$(find wasm-fauxrpc -type f -newer static/fauxrpc.wasm)" ]; then \
		echo "WASM is up to date."; \
	else \
		echo "Building WASM..."; \
		rm -f static/fauxrpc.wasm; \
		GOOS=js GOARCH=wasm go build -C wasm-fauxrpc -ldflags="-s -w" -o ../static/fauxrpc.wasm main.go; \
		wasm-opt -Oz --all-features static/fauxrpc.wasm -o static/fauxrpc.wasm; \
		gzip -9 -f static/fauxrpc.wasm; \
		mv static/fauxrpc.wasm.gz static/fauxrpc.wasm; \
	fi

fmt:
	buf format -w

lint:
	buf lint
