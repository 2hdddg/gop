package tree

import (
	"fmt"
	"github.com/2hdddg/gop/parser"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func assertPath(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("Expected path %v but was %v", expected, actual)
	}
}

func assertName(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("Expected name %v but was %v", expected, actual)
	}
}

func assertEmptyListOfPackages(t *testing.T, packs []*Package) {
	if packs == nil || len(packs) != 0 {
		t.Errorf("Expected non nil empty list of packages")
	}
}

func assertEmptyListOfFiles(t *testing.T, files []*File) {
	if files == nil || len(files) != 0 {
		t.Errorf("Expected non nil empty list of files")
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Exected no error but was %v", err)
	}
}

func onParsedFake(t *Tree, p *Package) {
}

type FakeParser struct {
}

func (p *FakeParser) Parse(path string) (*parser.Symbols, error) {
	return parser.NewSymbols(), nil
}

type FakeFile struct {
	name string
	path string
}

type FakeDir struct {
	name  string
	files []*FakeFile
	dirs  []*FakeDir
}

func (f *FakeFile) Name() string {
	return f.name
}

func (f *FakeFile) Size() int64 {
	return 0
}

func (f *FakeFile) Mode() os.FileMode {
	return 0
}

func (f *FakeFile) ModTime() time.Time {
	return time.Now()
}

func (f *FakeFile) IsDir() bool {
	return false
}

func (f *FakeFile) Sys() interface{} {
	return nil
}

func (f *FakeDir) Name() string {
	return f.name
}

func (f *FakeDir) Size() int64 {
	return 0
}

func (f *FakeDir) Mode() os.FileMode {
	return 0
}

func (f *FakeDir) ModTime() time.Time {
	return time.Now()
}

func (f *FakeDir) IsDir() bool {
	return true
}

func (f *FakeDir) Sys() interface{} {
	return nil
}

func (f *FakeDir) findDir(name string) *FakeDir {
	for _, d := range f.dirs {
		if d.name == name {
			return d
		}
	}
	return nil
}

func (root *FakeDir) ReadDirectory(dirpath string) (fis []os.FileInfo, err error) {
	log.Printf("In readdir %v", dirpath)
	parts := strings.Split(filepath.ToSlash(dirpath), "/")
	if len(parts) == 0 {
		err = fmt.Errorf("Empty path")
		return
	}

	filtered := parts[:0]
	for _, p := range parts {
		if p != "" {
			filtered = append(filtered, p)
		}
	}

	curr := root
	for _, part := range filtered {
		curr = curr.findDir(part)
		if curr == nil {
			err = fmt.Errorf("Path not found: %v", part)
			return
		}
	}
	for _, dir := range curr.dirs {
		fis = append(fis, dir)
	}
	for _, file := range curr.files {
		fis = append(fis, file)
	}
	return
}
