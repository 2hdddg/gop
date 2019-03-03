package server

import (
	"log"
	"os"
	"path"
)

type root struct {
	path string
}

type config struct {
	system    *root
	workspace *root
}

func newRoot(p string) *root {
	if p == "" {
		return nil
	}

	return &root{path: p}
}

func newConfig() *config {
	p := os.Getenv("GOROOT")
	if p == "" {
		p = "/usr/lib/go"
	}
	p = path.Join(p, "src")
	system := newRoot(p)

	workspace := newRoot(os.Getenv("GOPATH"))

	return &config{
		workspace: workspace,
		system:    system,
	}
}

func (c *config) valid() bool {
	if c.system == nil {
		return false
	}

	// Should be a valid GOROOT directory. Doesn't make sense
	// when this isn't valid..
	info, err := os.Stat(c.system.path)
	if err != nil || !info.IsDir() {
		log.Printf("Invalid GOROOT: %s", c.system.path)
		return false
	}

	// Optional if GOPATH is set or not, could be added later
	// or workspace will not be searchable.
	if c.workspace != nil {
		info, err = os.Stat(c.workspace.path)
		if err != nil || !info.IsDir() {
			log.Printf("Invalid GOPATH: %s", c.workspace.path)
			return false
		}
	}

	return true
}
