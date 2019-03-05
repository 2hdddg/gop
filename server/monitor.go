package server

import (
	"github.com/2hdddg/gop/shared"
	"io/ioutil"
	"log"
	"path"
	"strings"
	"sync"
)

type monitor struct {
	fileChan      chan *file
	systemPath    string
	workspacePath string
}

func newMonitor(config *shared.Config, fileChan chan *file) *monitor {
	return &monitor{fileChan: fileChan, systemPath: config.SystemPath}
}

func (m *monitor) _parseFiles(root, p string, files []string) {
	var waitGroup sync.WaitGroup

	for _, tmp := range files {
		waitGroup.Add(1)
		f := tmp
		go func() {
			defer waitGroup.Done()
			parsed, err := parseFile(p, path.Join(root, p, f))
			if err == nil {
				m.fileChan <- parsed
			}
		}()
	}
	waitGroup.Wait()
}

func _probe(path string) (dirs, files []string) {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		log.Printf("Error analyzing package at %v: %v\n", path, err)
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

func (m *monitor) _analyzePackage(root, packPath string) {
	fullPath := path.Join(root, packPath)
	log.Printf("Analyzing package at %s", fullPath)

	dirs, files := _probe(fullPath)
	m._parseFiles(root, packPath, files)

	for _, dir := range dirs {
		m._analyzePackage(root, path.Join(packPath, dir))
	}
}

func (m *monitor) _analyzeRoots() {
	m._analyzePackage(m.systemPath, "")
}

// Should not block too long
func (m *monitor) start() {
	go m._analyzeRoots()
}
