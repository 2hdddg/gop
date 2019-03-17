package tree

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/2hdddg/gop/parser"
)

type OnParsed func(t *Tree, p *Package)
type ReadDir func(dirname string) ([]os.FileInfo, error)

type Builder struct {
	onParsed OnParsed
	readDir  ReadDir
	tree     *Tree
	parse    Parse
}

func NewBuilder(onParsed OnParsed, rootPath string) (*Builder, error) {
	tree, err := NewTree(rootPath)
	if err != nil {
		return nil, err
	}

	return &Builder{
		onParsed: onParsed,
		readDir:  ioutil.ReadDir,
		tree:     tree,
		parse:    parser.Parse,
	}, nil
}

func (b *Builder) Build() (*Tree, error) {
	dirs, _, err := b.probe(b.tree.Path)
	if err != nil {
		return nil, err
	}

	for _, dir := range dirs {
		err = b.pack(dir, filepath.Join(b.tree.Path, dir))
	}
	return b.tree, nil
}

func (b *Builder) pack(name, path string) (err error) {
	dirs, files, err := b.probe(path)
	if err != nil {
		return
	}

	p := b.tree.AddPackage(name)

	for _, file := range files {
		p.AddFile(file, parser.Parse)
	}
	// Notify about completed package
	b.onParsed(b.tree, p)
	log.Printf("Parsed package %v", name)

	for _, dir := range dirs {
		b.pack(dir, filepath.Join(p.Path, dir))
	}

	return
}

func (b *Builder) probe(dir string) (dirs, files []string, err error) {
	fis, err := b.readDir(dir)
	if err != nil {
		return
	}

	for _, fi := range fis {
		name := fi.Name()
		mode := fi.Mode()

		// Ignore hidden dirs and files
		if strings.HasPrefix(name, ".") {
			continue
		}

		if fi.IsDir() {
			// No special check for dirs
			dirs = append(dirs, name)
			continue
		}

		// Checks for files start here

		if !mode.IsRegular() {
			continue
		}

		if !strings.HasSuffix(name, ".go") {
			continue
		}

		files = append(files, name)
	}

	return
}
