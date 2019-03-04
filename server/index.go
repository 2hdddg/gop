package server

import ()

type index struct {
	packs map[string]*pack
	funcs map[string][]*pack
}

func (i *index) packByName(name string) *Answer {
	var locations []Location

	p := i.packs[name]
	if p != nil {
		locations = append(locations, Location{Path: p.path})
	}
	return &Answer{Locations: locations}
}

func (i *index) funcByQuery(query *Query) *Answer {
	var locations []Location
	checkImported := len(query.Packages) > 0

	packs := i.funcs[query.Name]
	for _, hit := range packs {
		found := hit
		if checkImported {
			found = nil
			for _, imported := range query.Packages {
				if hit.name == imported {
					found = hit
					break
				}
			}
		}

		if found != nil {
			l := found.findFunc(query.Name)
			if l != nil {
				locations = append(locations, *l)
			}
		}
	}
	return &Answer{Locations: locations}
}
