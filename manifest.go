package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"tasc/patcher"

	"gopkg.in/yaml.v2"
)

// LoadError is for when the manifest file won't load.
type LoadError struct {
	msg string
}

// Error returns the load error message.
func (e LoadError) Error() string {
	return e.msg
}

// ParseError is for when the manifest won't parse.
type ParseError struct {
	msg string
}

// Error returns the parse error message.
func (e ParseError) Error() string {
	return e.msg
}

// The Manifest is the structural representation of the manifest.
type Manifest struct {
	Projects []*Project
	Patches  []*patcher.Patch
}

// BlockingProjects returns slices of the blocking and non-blocking projects.
func (m *Manifest) BlockingProjects() (SortProject, SortProject) {
	var s SortProject
	var a SortProject

	for _, project := range m.Projects {
		if project.Blocking {
			s = append(s, project)
		} else {
			a = append(a, project)
		}
	}

	return s, a
}

// UnmarshalYAML is an implementation of the YAML Unmarshaler interface so we
// can have better control over how a Manifest us created from YAML.
func (m *Manifest) UnmarshalYAML(unmarshal func(interface{}) error) error {
	f := make(map[string][]map[string]interface{})

	// First, lets get the original unmarshalled value
	if err := unmarshal(f); err != nil {
		return err
	}

	// Projects
	for _, pr := range f["projects"] {
		project := NewProjectFromMap(pr)
		m.Projects = append(m.Projects, project)
	}

	// Patches
	for _, pa := range f["patches"] {
		patch := patcher.NewPatchFromMap(pa)
		m.Patches = append(m.Patches, patch)
	}

	return nil
}

// Load a manifest from a yaml filename and a map of params.
func (m *Manifest) Load(filename string, params map[string]string) error {
	// Read the manifest YAML file.
	mb, err := ioutil.ReadFile(filename)
	if err != nil {
		return LoadError{"Error loading"}
	}

	// Replace any instances of {param} with the value from params in the
	// string represenation of the manifest (ms).
	ms := string(mb)
	for param, value := range params {
		ms = strings.Replace(ms, fmt.Sprintf("{%s}", param), value, -1)
	}

	// Parse the manifest YAML.
	err = yaml.Unmarshal([]byte(ms), &m)
	if err != nil {
		return ParseError{"Error parsing"}
	}

	return nil
}
