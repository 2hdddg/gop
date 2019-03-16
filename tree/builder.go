package tree

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/2hdddg/gop/parser"
)

type OnParsed func(t *Tree, p *Package)

type Builder struct {
	onParsed OnParsed
	readDir  func(dirname string) ([]os.FileInfo, error)
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

func (b *Builder) Build() (t *Tree, err error) {
	dirs, _, err := b.probe(b.tree.Path)
	if err != nil {
		return
	}

	for _, dir := range dirs {
		err = b.pack(dir)
	}
	return
}

func (b *Builder) pack(dir string) (err error) {
	dirs, files, err := b.probe(dir)
	if err != nil {
		return
	}

	p := b.tree.AddPackage(dir)

	for _, file := range files {
		p.AddFile(file, parser.Parse)
	}
	// Notify about completed package
	b.onParsed(b.tree, p)

	for _, dir := range dirs {
		b.pack(dir)
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

		if mode.IsDir() {
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
