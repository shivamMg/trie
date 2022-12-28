.PHONY: build-wasm
build-wasm:
	cd wasm && GOOS=js GOARCH=wasm go build -o ../site/main.wasm

.PHONY: server
server:
	cd server && go run .
