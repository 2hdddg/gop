package server

import ()

type index struct {
	packs map[string]*pack
	funcs map[string][]*pack
}

func (i *index) packByName(name string) *Answer {
	var locations []FileLocation

	p := i.packs[name]
	if p != nil {
		locations = append(locations, FileLocation{FilePath: name})
	}
	return &Answer{Locations: locations}
}

func (i *index) funcByName(name string) *Answer {
	var locations []FileLocation

	packs := i.funcs[name]
	for _, p := range packs {
		l := p.findFunc(name)
		if l != nil {
			locations = append(locations, *l)
		}
	}
	return &Answer{Locations: locations}
}
