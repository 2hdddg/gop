package shared

import (
	"log"
	"os"
	"path"
)

type Config struct {
	SystemPath    string
	WorkspacePath string
}

func NewConfig() *Config {
	p := os.Getenv("GOROOT")
	if p == "" {
		p = "/usr/lib/go"
	}
	p = path.Join(p, "src")
	system := p

	workspace := os.Getenv("GOPATH")

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
