package tree

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type DirectoryReader interface {
	ReadDirectory(path string) ([]os.FileInfo, error)
}

type Progress interface {
	OnPackageParsed(t *Tree, p *Package)
	OnTreeParsed(t *Tree)
}

type Builder struct {
	Progress Progress
	reader   DirectoryReader
	tree     *Tree
	parser   Parser
}

func NewBuilder(rootPath string) (*Builder, error) {
	tree, err := NewTree(rootPath)
	if err != nil {
		return nil, err
	}

	builder := &Builder{
		tree:   tree,
		parser: tree,
	}
	builder.reader = builder

	return builder, nil
}

// Implement DirecoryReader
func (b *Builder) ReadDirectory(path string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(path)
}

func (b *Builder) Build() (*Tree, error) {
	dirs, _, err := b.probe(b.tree.Path)
	if err != nil {
		return nil, err
	}

	for _, dir := range dirs {
		err = b.pack(dir, filepath.Join(b.tree.Path, dir))
	}

	if b.Progress != nil {
		b.Progress.OnTreeParsed(b.tree)
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
		p.AddFile(file, b.parser)
	}
	// Notify about completed package
	if b.Progress != nil {
		b.Progress.OnPackageParsed(b.tree, p)
	}

	for _, dir := range dirs {
		fullName := name + "/" + dir
		b.pack(fullName, filepath.Join(p.Path, dir))
	}

	return
}

func (b *Builder) probe(dir string) (dirs, files []string, err error) {
	fis, err := b.reader.ReadDirectory(dir)
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
