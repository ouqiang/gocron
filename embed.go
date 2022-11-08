package gocron

import (
	"embed"
	"io/fs"
)

//go:embed web/vue
var embedFrontend embed.FS

func GetFrontendAssets() (fs.FS, error) {
	return fs.Sub(embedFrontend, "web/vue/dist")
}
