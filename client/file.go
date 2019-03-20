package client

import (
	"github.com/2hdddg/gop/config"
	"path/filepath"
	"strings"
)

func parseFilePackage(config *config.Config,
	path string) (pack string, e error) {
	p, _ := filepath.Abs(path)
	p = filepath.Dir(p)
	p = strings.TrimPrefix(p, config.WorkspacePath)
	pack = strings.TrimPrefix(p, "/")
	return
}
