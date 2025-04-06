package web

import (
	"embed"
	"io/fs"
	"log/slog"
	"os"
)

//go:embed build
var webFS embed.FS

func GetWebFS() fs.FS {
	embedRoot, err := fs.Sub(webFS, "build")
	if err != nil {
		slog.Error("Unable to get root for web ui", "error", err)
		os.Exit(1)
	}
	return embedRoot
}
