package server

import (
	"log"
	"os"
	"path"
)

type config struct {
	goRoot string
	goPath string
}

func newConfig() *config {
	goRoot := os.Getenv("GOROOT")
	if goRoot == "" {
		goRoot = "/usr/lib/go"
	}
	goRoot = path.Join(goRoot, "src")

	return &config{
		goPath: os.Getenv("GOPATH"),
		goRoot: goRoot,
	}
}

func (c *config) valid() bool {
	// Should be a valid GOROOT directory. Doesn't make sense
	// when this isn't valid..
	info, err := os.Stat(c.goRoot)
	if err != nil || !info.IsDir() {
		log.Printf("Invalid GOROOT: %s", c.goRoot)
		return false
	}

	// Optional if GOPATH is set or not, could be added later
	// or workspace will not be searchable.
	if c.goPath != "" {
		info, err = os.Stat(c.goPath)
		if err != nil || !info.IsDir() {
			log.Printf("Invalid GOPATH: %s", c.goPath)
			return false
		}
	}

	return true
}
