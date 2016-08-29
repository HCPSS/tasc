package fetcher

// Fetcher is an interface for types that fetch source code.
type Fetcher interface {
	GetSource() string
	GetDestination() string
	Fetch(baseDir string) error
}
