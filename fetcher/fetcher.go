package fetcher

import "sync"

// Fetcher is an interface for types that fetch source code.
type Fetcher interface {
	GetSource() string
	GetDestination() string
	Fetch(baseDir string, wg *sync.WaitGroup) error
}
