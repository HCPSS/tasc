package main

import (
	"sort"
	"sync"
	"tasc/patcher"
)

// Tasc is the main structure for the application. It is responsible for
// fetching and patching projects.
type Tasc struct {
	manifest    Manifest
	destination string
}

// Fetch fetches the project and updates the progress.
func Fetch(proj *Project, dest string, prog *Progress, wg *sync.WaitGroup) {
	prog.Add(proj, StateProcessing).Report()
	if err := proj.Fetcher.Fetch(dest, wg); err != nil {
		prog.Add(proj, StateFailed).Report()
	} else {
		prog.Add(proj, StateSuccess).Report()
	}
}

// Assemble the source code for the project.
func (t *Tasc) Assemble(c chan string) {
	progress := NewProgress(c)
	progress.QueueProjects(t.manifest.Projects)

	wg := &sync.WaitGroup{}
	wg.Add(len(t.manifest.Projects))

	// Synchronous and asynchronous projects.
	sProjs, aProjs := t.manifest.BlockingProjects()

	// First lets work through the synchronous projects.
	sort.Sort(sProjs)
	for _, sProj := range sProjs {
		Fetch(sProj, t.destination, progress, wg)
	}

	// Now the asyncronous ones.
	sort.Sort(aProjs)
	for _, aProj := range aProjs {
		go func(p *Project) {
			Fetch(p, t.destination, progress, wg)
		}(aProj)
	}

	go func() {
		wg.Wait()
		close(c)
	}()
}

// Patch performs the patches.
func (t *Tasc) Patch() patcher.PatchResults {
	var results patcher.PatchResults

	for _, patch := range t.manifest.Patches {
		result := patch.Patcher.Patch()
		results = append(results, result)
	}

	return results
}
