wasm:
	GOOS=js GOARCH=wasm go build -mod vendor -o www/wasm/update_exif.wasm cmd/update-exif-wasm/main.go

cli:
	go build -mod vendor -o bin/update-exif cmd/update-exif/main.go
	go build -mod vendor -o bin/server cmd/update-exif-server/main.go
