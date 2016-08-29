package fetcher

import (
	"io"
	"os"
	"path/filepath"
)

// LocalFetcher fetches local files
type LocalFetcher struct {
	source, destination string
}

// NewLocalFetcher gets a new LocalFetcher
func NewLocalFetcher(source, destination string) *LocalFetcher {
	lf := new(LocalFetcher)

	lf.source = source
	lf.destination = destination

	return lf
}

// GetSource returns the location of the source code and is required by the
//Fetcher interface.
func (lf *LocalFetcher) GetSource() string {
	return lf.source
}

// GetDestination returns the location of the destination and is required by the
//Fetcher interface.
func (lf *LocalFetcher) GetDestination() string {
	return lf.destination
}

// Fetch fetches the source code and is required by the Fetcher interface.
func (lf *LocalFetcher) Fetch(baseDir string) error {
	s, err := os.Open(lf.source)
	if err != nil {
		return err
	}
	defer s.Close()

	info, err := s.Stat()
	if err != nil {
		return err
	}

	switch mode := info.Mode(); {
	case mode.IsDir():
		err = os.MkdirAll(lf.destination, 0755)
		if err != nil {
			return err
		}

		CopyDir(lf.source, lf.destination)
	case mode.IsRegular():
		err = os.MkdirAll(filepath.Dir(lf.destination), 0755)
		if err != nil {
			return err
		}

		CopyFile(lf.source, lf.destination)
	}

	return nil
}

// CopyFile copys a file.
func CopyFile(source, destination string) (err error) {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err == nil {
		sourceInfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(destination, sourceInfo.Mode())
		}
	}

	return
}

// CopyDir copys a directory.
func CopyDir(source, destination string) (err error) {
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	err = os.MkdirAll(destination, sourceInfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)
	objects, err := directory.Readdir(-1)
	for _, obj := range objects {
		sourceFilePointer := filepath.Join(source, obj.Name())
		destFilePointer := filepath.Join(destination, obj.Name())

		if obj.IsDir() {
			// Create subdirectories
			err = CopyDir(sourceFilePointer, destFilePointer)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(sourceFilePointer, destFilePointer)
			if err != nil {
				return err
			}
		}
	}

	return
}
