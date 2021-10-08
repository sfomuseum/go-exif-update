GOROOT=$(shell go env GOROOT)

tag-data:
	curl -o tags/tags_data.go 

rebuild-wasm:
	@make wasmjs
	@make wasm

wasm:
	GOOS=js GOARCH=wasm go build -mod vendor -o www/wasm/update_exif.wasm cmd/update-exif-wasm/main.go
	GOOS=js GOARCH=wasm go build -mod vendor -o www/wasm/supported_tags.wasm cmd/tags-supported-wasm/main.go

wasmjs:
	cp "$(GOROOT)/misc/wasm/wasm_exec.js" www/javascript/

cli:
	@make wasm
	go build -mod vendor -o bin/tags-is-supported cmd/tags-is-supported/main.go
	go build -mod vendor -o bin/tags-supported cmd/tags-supported/main.go
	go build -mod vendor -o bin/update-exif cmd/update-exif/main.go
	go build -mod vendor -o bin/update-exif-wasm-demo cmd/update-exif-wasm-demo/main.go

debug:
	go run -mod vendor cmd/update-exif-wasm-demo/main.go

lambda:
	@make lambda-server

lambda-server:
	if test -f main; then rm -f main; fi
	if test -f server.zip; then rm -f server.zip; fi
	GOOS=linux go build -mod vendor -o main cmd/update-exif-server/main.go
	zip server.zip main
	rm -f main
