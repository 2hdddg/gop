package server

import (
	"path"
)

type index struct {
	root  string
	packs map[string]*pack
	funcs map[string][]*pack
}

func (i *index) packByName(name string) *Answer {
	var locations []Location

	p := i.packs[name]
	if p != nil {
		location := Location{Path: path.Join(i.root, p.name)}
		locations = append(locations, location)
	}
	return &Answer{Locations: locations}
}

func (i *index) funcByQuery(query *Query) *Answer {
	var locations []Location
	checkImported := len(query.Packages) > 0

	packs := i.funcs[query.Name]
	for _, candidate := range packs {
		match := candidate
		if checkImported {
			match = nil
			for _, imported := range query.Packages {
				if candidate.name == imported {
					match = candidate
					break
				}
			}
		}

		if match != nil {
			l := match.findFunc(query.Name)
			if l != nil {
				locations = append(locations, *l)
			}
		}
	}
	return &Answer{Locations: locations}
}
