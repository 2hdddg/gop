// Tree is a representation of a parsed source code tree of
// Go code. The representation is used as input for building
// a search index. The tree should be able to handle updates
// to files and structure in runtime and rebuild search index
// when needed.

package server

import ()

type tree struct {
	root  string           // Root path: /usr/lib/go/src
	packs map[string]*pack // Set of packs used to build an index.
	dirty bool
}

func newTree(root string) *tree {
	return &tree{
		root:  root,
		packs: make(map[string]*pack),
		dirty: false,
	}
}

// Ensures that package referenced by parsed file exists in
// the set of packs.
func (t *tree) _ensurePack(f *file) *pack {
	pack, exists := t.packs[f.packName]
	if !exists {
		pack = newPack(f.packName)
		t.packs[f.packName] = pack
	}

	return pack
}

// Adds a pack to the set indirectly by adding a file that belongs to
// a package.
func (t *tree) addFile(f *file) {
	p := t._ensurePack(f)
	p.mergeFile(f)
	t.dirty = true
}

func (t *tree) buildIndex() *index {
	packs := make(map[string]*pack)
	funcs := make(map[string][]*pack)

	for _, p := range t.packs {
		// Important to copy package since it will be sent
		// to channel handling search while it might be modified
		// by channel that indexes packages.
		pcopy := *p
		packs[p.name] = &pcopy
		for n, _ := range pcopy.funcs {
			funcs[n] = append(funcs[n], &pcopy)
		}
	}
	t.dirty = false
	return &index{root: t.root, packs: packs, funcs: funcs}
}
