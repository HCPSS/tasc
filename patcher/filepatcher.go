package patcher

import (
	"fmt"
	"os/exec"
)

// FilePatcher patches based on patch files.
type FilePatcher struct {
	Source      string
	Destination string
}

// NewFilePatcher returns a new FilePatcher.
func NewFilePatcher(source, destination string) *FilePatcher {
	fp := new(FilePatcher)

	fp.Source = source
	fp.Destination = destination

	return fp
}

// Patch performs the patch. Needed to satisfy Patcher interface.
func (p *FilePatcher) Patch() *PatchResult {
	err := exec.Command(
		"patch",
		fmt.Sprintf("--input=%s", p.Source),
		p.Destination,
	).Run()

	return &PatchResult{Error: err, Patcher: p}
}

// GetSource sets the source patch file. Needed to satisfy Patcher interface.
func (p *FilePatcher) GetSource() string {
	return p.Source
}

// GetDestination sets the destination where the patch should be applied. Needed
// to satisfy Patcher interface.
func (p *FilePatcher) GetDestination() string {
	return p.Destination
}

// SetSource sets the source.
func (p *FilePatcher) SetSource(source string) {
	p.Source = source
}

// SetDestination sets the destination.
func (p *FilePatcher) SetDestination(destination string) {
	p.Destination = destination
}
