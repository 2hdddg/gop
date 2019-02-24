package server

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"sync"
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

func (r *root) build(f chan *file) {
	log.Printf("Building root at %s", r.path)
	buildDir(f, r.path)
}

func parseFiles(fileChan chan *file, p string, files []string) {
	var waitGroup sync.WaitGroup

	for _, tmp := range files {
		waitGroup.Add(1)
		f := tmp
		go func() {
			defer waitGroup.Done()
			parsed, err := parseFile(path.Join(p, f))
			if err == nil {
				fileChan <- parsed
			}
		}()
	}
	waitGroup.Wait()
}
func buildDir(f chan *file, p string) {
	log.Printf("Analyzing package at %s", p)

	dirs, files := probe(p)
	parseFiles(f, p, files)

	for _, dir := range dirs {
		buildDir(f, path.Join(p, dir))
	}
}

func probe(path string) (dirs, files []string) {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		log.Printf("Error indexing package at %v: %v\n", path, err)
		return nil, nil
	}

	for _, i := range entries {
		mode := i.Mode()
		name := i.Name()
		if mode.IsDir() {
			dirs = append(dirs, name)
		} else if mode.IsRegular() && strings.LastIndex(name, ".go") > 0 {
			files = append(files, name)
		}
	}
	return dirs, files
}
