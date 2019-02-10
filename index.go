package main

import (
	"fmt"
	"net/http"
)

type Index struct {
	packages map[string]Package
}

type Answer interface {
	Response(w http.ResponseWriter)
}

type PackagesAnswer struct {
	Packages []string
}

func (a *PackagesAnswer) Response(w http.ResponseWriter) {
	fmt.Fprintf(w, "%v", *a)
}

type LocationsAnswer struct {
	Locations []FileLocation
}

func (a *LocationsAnswer) Response(w http.ResponseWriter) {
	fmt.Fprintf(w, "%v", *a)
}

type Query interface {
	Process(i *Index) Answer
}

type PackagesQuery struct {
}

func (q *PackagesQuery) Process(i *Index) Answer {
	var packages []string
	for _, p := range i.packages {
		packages = append(packages, p.Name)
	}
	return &PackagesAnswer{Packages: packages}
}

type DefinitionQuery struct {
	name string
}

func (q *DefinitionQuery) Process(i *Index) Answer {
	var locations []FileLocation
	for _, p := range i.packages {
		l := p.Find(q.name)
		if l != nil {
			locations = append(locations, *l)
		}
	}
	return &LocationsAnswer{Locations: locations}
}
