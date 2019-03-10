package server

import (
	"io/ioutil"
	"log"
	"path"
	"strings"
)

type traverser struct {
	root     string
	fileChan chan *file
}

func newTraverser(root string, fileChan chan *file) *traverser {
	return &traverser{
		root:     root,
		fileChan: fileChan,
	}
}

func (t *traverser) traverse() {
	log.Printf("Traversing %v", t.root)
	dirs, _ := _probe(t.root)
	// Each subdirectory to root is a package
	for _, pack := range dirs {
		t._analyze(pack)
	}
}

func (t *traverser) _analyze(pack string) {
	packPath := path.Join(t.root, pack)

	log.Printf("Analyzing %v", pack)

	dirs, files := _probe(packPath)
	t._parse(pack, files)

	for _, dir := range dirs {
		// Each subdirectory to this package is another package
		t._analyze(path.Join(pack, dir))
	}
}

func (t *traverser) _parse(pack string, files []string) {
	for _, filename := range files {
		filePath := path.Join(t.root, pack, filename)
		f := newFile(pack, filePath)
		f.parse()
		if f.valid {
			log.Printf("Sending %v", filePath)
			t.fileChan <- f
		}
	}
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
