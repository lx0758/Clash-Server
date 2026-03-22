package res

import (
	"embed"
	"io/fs"
)

var (
	//go:embed all:web/*
	webEmbedFS     embed.FS
	WebFS, _       = fs.Sub(webEmbedFS, "web")
	WebAssetsFS, _ = fs.Sub(WebFS, "assets")
)
