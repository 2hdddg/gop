package server

import (
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

func newMonitor(config *config, fileChan chan *file) *monitor {
	return &monitor{fileChan: fileChan, systemPath: config.system.path}
}

func (m *monitor) _parseFiles(p string, files []string) {
	var waitGroup sync.WaitGroup

	for _, tmp := range files {
		waitGroup.Add(1)
		f := tmp
		go func() {
			defer waitGroup.Done()
			parsed, err := parseFile(path.Join(p, f))
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

func (m *monitor) _analyzePackage(p string) {
	log.Printf("Analyzing package at %s", p)

	dirs, files := _probe(p)
	m._parseFiles(p, files)

	for _, dir := range dirs {
		m._analyzePackage(path.Join(p, dir))
	}
}

func (m *monitor) _analyzeRoots() {
	m._analyzePackage(m.systemPath)
}

// Should not block too long
func (m *monitor) start() {
	go m._analyzeRoots()
}
