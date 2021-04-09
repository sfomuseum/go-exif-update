// package update provides methods for updating EXIF data in JPEG files.
// This is a thin wrapper around code in dsoprea's go-exif and go-jpeg-image-structure packages
// and includes command-line tools for updating the EXIF data JPEG files using key-value parameters
// as well as a WebAssembly (wasm) binary for updating EXIF data in JavaScript (or other languages
// that support wasm binaries).
//
// Importantly not all EXIF tags are supported yet.
package update
