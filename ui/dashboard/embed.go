package dashboard

import (
	"embed"
	"io/fs"
)

//go:embed build
var assets embed.FS

// FS contains the new dashboard assets.
var FS, _ = fs.Sub(assets, "build")
