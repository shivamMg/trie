.PHONY: bindata
bindata:
	cd wasm && go-bindata -nometadata -nocompress words.txt

.PHONY: build-wasm
build-wasm:
	cd wasm && GOOS=js GOARCH=wasm go build -o ../site/main.wasm

.PHONY: server
server:
	cd server && go run .
