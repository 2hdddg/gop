package config

import (
	"log"
	"os"
	"path"
	"runtime"
)

type Config struct {
	SystemPath    string
	WorkspacePath string
}

func NewConfig() *Config {
	system := path.Join(runtime.GOROOT(), "src")
	// GOPATH is a list of paths (colon-separated on Unix)
	// TODO: Handle multiple paths!
	workspace := os.Getenv("GOPATH")
	if workspace != "" {
		workspace = path.Join(workspace, "src")
	}

	return &Config{
		WorkspacePath: workspace,
		SystemPath:    system,
	}
}

func (c *Config) Valid() bool {
	if c.SystemPath == "" {
		return false
	}

	// Should be a valid GOROOT directory. Doesn't make sense
	// when this isn't valid..
	info, err := os.Stat(c.SystemPath)
	if err != nil || !info.IsDir() {
		log.Printf("Invalid GOROOT: %s", c.SystemPath)
		return false
	}

	// Optional if GOPATH is set or not, could be added later
	// or workspace will not be searchable.
	if c.WorkspacePath != "" {
		info, err = os.Stat(c.WorkspacePath)
		if err != nil || !info.IsDir() {
			log.Printf("Invalid GOPATH: %s", c.WorkspacePath)
			return false
		}
	}

	return true
}
