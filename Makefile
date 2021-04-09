tag-data:
	curl -o tags/tags_data.go https://raw.githubusercontent.com/dsoprea/go-exif/de2141190595193aa097a2bf3205ba0cf76dc14b/tags_data.go

wasm:
	GOOS=js GOARCH=wasm go build -mod vendor -o www/wasm/update_exif.wasm cmd/update-exif-wasm/main.go
	GOOS=js GOARCH=wasm go build -mod vendor -o www/wasm/supported_tags.wasm cmd/tags-supported-wasm/main.go

cli:
	go build -mod vendor -o bin/update-exif cmd/update-exif/main.go
	go build -mod vendor -o bin/server cmd/update-exif-server/main.go

debug:
	go run -mod vendor cmd/update-exif-server/main.go

lambda:
	@make lambda-server

lambda-server:
	if test -f main; then rm -f main; fi
	if test -f server.zip; then rm -f server.zip; fi
	GOOS=linux go build -mod vendor -o main cmd/update-exif-server/main.go
	zip server.zip main
	rm -f main
