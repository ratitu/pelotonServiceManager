package web

import "embed"

// Files contains embedded web UI assets.
//
//go:embed index.html style.css
var Files embed.FS
