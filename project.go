package main

import (
	"reflect"
	"strings"
	"tasc/fetcher"
)

// Project is the representation of an individual project.
type Project struct {
	Name     string
	Fetcher  fetcher.Fetcher
	Blocking bool
	Sticky   bool
}

// SortProject is a sortable list of Project.
type SortProject []*Project

// The Len, Swap and Less methods are needed to satisfy the sort.Interface.
func (c SortProject) Len() int      { return len(c) }
func (c SortProject) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c SortProject) Less(i, j int) bool {
	// Is c[i] less than c[j]?
	if (c[i].Sticky && c[j].Sticky) || (!c[i].Sticky && !c[j].Sticky) {
		// Either both c[i] and c[j] are sticky, or neither are. In either case
		// they can be compared normally.
		return c[i].Name < c[j].Name
	}

	return c[j].Sticky
}

// NewProjectFromMap creates a new Project from a map.
func NewProjectFromMap(mp map[string]interface{}) *Project {
	project := new(Project)

	source, _ := mp["source"].(string)
	destination, _ := mp["destination"].(string)

	// Name
	tokens := strings.Split(source, "/")
	project.Name = tokens[len(tokens)-1]

	// Fetcher
	switch mp["provider"] {
	case "git":
		rename, _ := mp["rename"].(string)
		version, _ := mp["version"].(string)

		project.Fetcher = fetcher.NewGitFetcher(
			source, destination, rename, version,
		)
	case "zip":
		fallthrough
	default:
		project.Fetcher = fetcher.NewArchiveFetcher(
			source, destination,
		)
	}

	// Tags
	tags, ok := mp["tags"]
	if ok {
		switch reflect.TypeOf(tags).Kind() {
		case reflect.Slice:
			t := reflect.ValueOf(tags)

			for i := 0; i < t.Len(); i++ {
				switch t.Index(i).Interface() {
				case "blocking":
					project.Blocking = true
				case "sticky":
					project.Sticky = true
				}
			}
		}
	}

	return project
}
