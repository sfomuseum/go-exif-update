package static

import (
	"embed"
)

//go:embed css/* javascript/*
var FS embed.FS
