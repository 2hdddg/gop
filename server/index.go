package server

import ()

type index struct {
	packs map[string]pack
}

func (i *index) allPackages() *PackagesAnswer {
	var packs []string
	for _, p := range i.packs {
		packs = append(packs, p.name)
	}
	return &PackagesAnswer{Packages: packs}
}

func (i *index) funcDefinition(name string) *LocationsAnswer {
	var locations []FileLocation
	for _, p := range i.packs {
		l := p.find(name)
		if l != nil {
			locations = append(locations, *l)
		}
	}
	return &LocationsAnswer{Locations: locations}
}
