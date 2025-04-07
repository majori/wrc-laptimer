package web

import (
	"embed"
	"io/fs"
)

//go:embed index.html styles.css
var webFS embed.FS

func GetWebFS() fs.FS {
	return webFS
}
