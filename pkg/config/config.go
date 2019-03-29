package config

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

type Config struct {
	SystemPath    string
	WorkspacePath string
}

func evalPath(in string) (p string, err error) {
	// Path should be absolute to simplify comparisons
	if p, err = filepath.Abs(in); err != nil {
		return
	}
	// Symbolic links are followed until real target found, we
	// should use that path instead!
	if p, err = filepath.EvalSymlinks(p); err != nil {
		return
	}
	// Should be a valid path
	var i os.FileInfo
	if i, err = os.Stat(p); err != nil {
		return
	}
	// Should be a directory
	if !i.IsDir() {
		err = errors.New(fmt.Sprintf("%v is not a directory", p))
		return
	}

	// If there is a subdir named src, use that.
	// Note that the src dir also could be a symlink.
	psrc := path.Join(p, "src")
	if i, err = os.Stat(psrc); err != nil {
		err = nil
		return
	}
	if i.IsDir() {
		p, _ = filepath.EvalSymlinks(psrc)
	}

	return
}

func NewConfig() *Config {
	system := runtime.GOROOT()
	system, _ = evalPath(system)

	// GOPATH is a list of paths (colon-separated on Unix)
	// TODO: Handle multiple paths!
	workspace := os.Getenv("GOPATH")
	if workspace != "" {
		workspace, _ = evalPath(workspace)
	}

	return &Config{
		WorkspacePath: workspace,
		SystemPath:    system,
	}
}

func (c *Config) PackageFromPath(path string) (string, bool) {
	prefixes := []string{c.WorkspacePath, c.SystemPath}

	path, _ = filepath.Abs(filepath.Dir(path))
	for _, prefix := range prefixes {
		if prefix != "" && strings.HasPrefix(path, prefix) {
			pack := strings.TrimPrefix(path, prefix)
			pack = strings.TrimPrefix(pack, "/")
			return pack, true
		}
	}
	return "", false
}

func (c *Config) Paths() []string {
	paths := []string{c.SystemPath}
	if c.WorkspacePath != "" {
		paths = append(paths, c.WorkspacePath)
	}
	return paths
}
