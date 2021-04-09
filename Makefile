wasm:
	GOOS=js GOARCH=wasm go build -mod vendor -o www/wasm/update_exif.wasm cmd/update-exif-wasm/main.go
