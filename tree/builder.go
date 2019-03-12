package tree

import (
	"io/ioutil"
	"os"
	"strings"
)

type OnParsed func(t *Tree, p *Package)

type Builder struct {
	onParsed OnParsed
	readDir  func(dirname string) ([]os.FileInfo, error)
}

func NewBuilder(onParsed OnParsed) *Builder {
	return &Builder{
		onParsed: onParsed,
		readDir:  ioutil.ReadDir,
	}
}

func (b *Builder) Build(rootPath string) (t *Tree, err error) {
	t = NewTree(rootPath)

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
