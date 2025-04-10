package web

import (
	"embed"
	"io/fs"
)

//go:embed index.html styles.css index.js api.js bg.jpg
var webFS embed.FS

func GetWebFS() fs.FS {
	return webFS
}
