package fetcher

import (
	"path/filepath"
	"os/exec"
)

// NewSvnFetcher gets a new new SvnFetcher
func NewSvnFetcher(source, destination, rename, version string) *SvnFetcher {
	sf := new(SvnFetcher)

	sf.source = source
	sf.destination = destination
	sf.rename = rename
	sf.version = version

	return sf
}

// SvnFetcher fetches source code from git.
type SvnFetcher struct {
	rename, source, destination, version string
}

// GetSource returns the location of the source code and is required by the
// Fetcher interface.
func (sf *SvnFetcher) GetSource() string {
	return sf.source
}

// GetDestination returns the location of the destination and is required by the
// Fetcher interface.
func (sf *SvnFetcher) GetDestination() string {
	return sf.destination
}

// Fetch fetches the source code and is required by the Fetcher interface.
func (sf *SvnFetcher) Fetch(baseDir string) error {
	dest := filepath.Join(baseDir, sf.destination, sf.rename)
	cmd := exec.Command("svn", "co", sf.source, dest)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
