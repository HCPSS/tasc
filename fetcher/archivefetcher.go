package fetcher

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pivotal-golang/archiver/extractor"
)

// ArchiveFetcher fetches source code from a remote archive.
type ArchiveFetcher struct {
	source, destination string
}

// GetSource gets the path to the source and is required by the Fetcher
// interface.
func (af *ArchiveFetcher) GetSource() string {
	return af.source
}

// GetDestination gets the path to the destination and is required by the
// Fetcher interface.
func (af *ArchiveFetcher) GetDestination() string {
	return af.destination
}

func (af *ArchiveFetcher) downloadFile() (*os.File, error) {
	// Figure out the filename
	tokens := strings.Split(af.source, "/")
	fileName := tokens[len(tokens)-1]

	// Create the file that we will write to.
	tempFile, err := ioutil.TempFile("", fileName)
	if err != nil {
		return nil, err
	}
	defer tempFile.Close()

	// Download the remote file
	response, err := http.Get(af.source)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Write the downloaded file to the temp file
	_, err = io.Copy(tempFile, response.Body)
	if err != nil {
		return nil, err
	}

	return tempFile, err
}

// Fetch the source code. Required by the Fetcher interface.
func (af *ArchiveFetcher) Fetch(baseDir string) error {
	file, err := af.downloadFile()
	if err != nil {
		return err
	}

	dest := filepath.Join(baseDir, af.destination)

	extractor := extractor.NewDetectable()
	extractor.Extract(file.Name(), dest)

	// Clean up
	os.Remove(file.Name())

	return nil
}

// NewArchiveFetcher gets a new ArchiveFetcher.
func NewArchiveFetcher(source, destination string) *ArchiveFetcher {
	zf := new(ArchiveFetcher)

	zf.source = source
	zf.destination = destination

	return zf
}
