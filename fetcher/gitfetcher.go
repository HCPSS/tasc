package fetcher

import (
	"path/filepath"
	"sync"
	"time"

	"github.com/gogits/git-module"
)

// NewGitFetcher gets a new GitFetcher.
func NewGitFetcher(source, destination, rename, version string) *GitFetcher {
	gf := new(GitFetcher)

	gf.source = source
	gf.destination = destination
	gf.rename = rename
	gf.version = version

	return gf
}

// GitFetcher fetches source code from git.
type GitFetcher struct {
	rename, source, destination, version string
}

// GetSource returns the location of the source code and is required by the
//Fetcher interface.
func (gf *GitFetcher) GetSource() string {
	return gf.source
}

// GetDestination returns the location of the destination and is required by the
//Fetcher interface.
func (gf *GitFetcher) GetDestination() string {
	return gf.destination
}

// Fetch fetches the source code and is required by the Fetcher interface.
func (gf *GitFetcher) Fetch(baseDir string, wg *sync.WaitGroup) error {
	defer wg.Done()

	dest := filepath.Join(baseDir, gf.destination, gf.rename)
	cro := git.CloneRepoOptions{Timeout: time.Minute * 5}
	err := git.Clone(gf.source, dest, cro)
	if err != nil {
		return err
	}

	// Hopefully support for "checkout" will be added nativly:
	// https://github.com/gogits/git-module/pull/11
	_, err = git.NewCommand("checkout", gf.version).RunInDir(dest)
	if err != nil {
		panic(err)
	}

	return nil
}
