package files

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"
)

// Local is an implementation of Storage interface for storing files on disk locally.
type Local struct {
	maxFileSize int //max number of bytes for a file
	basePath    string
}

// NewLocal creates a new Local storage.
func NewLocal(basePath string) (*Local, error) {
	fmt.Println("basepath", basePath)
	p, err := filepath.Abs(basePath)
	fmt.Println("abs-basepath", p)
	if err != nil {
		return nil, err
	}
	return &Local{basePath: p}, nil
}

func (l *Local) Save(path string, contents io.Reader) error {
	fp := l.fullpath(path)
	fmt.Println("fullpath", fp)

	d := filepath.Dir(fp)
	fmt.Println("dirfp", d)
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return xerrors.Errorf("Unable to create directory: %w", err)
	}
	_, err = os.Stat(fp)
	if err == nil {
		err = os.Remove(fp)
		if err != nil {
			return xerrors.Errorf("Unable to remove existing file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return xerrors.Errorf("Unable to check if file exists: %w", err)
	}

	//create a new file at path
	f, err := os.Create(fp)
	if err != nil {
		return xerrors.Errorf("Unable to create file: %w", err)
	}
	defer f.Close()

	//write contents to file
	_, err = io.Copy(f, contents)
	if err != nil {
		return xerrors.Errorf("Unable to write to file: %w", err)
	}
	return nil
}

// Get file at given path and returna a reader
// the calling function is responsible for closing the reader
func (l *Local) Get(path string) (*os.File, error) {
	// get the full path for a file
	fp := l.fullpath(path)
	fmt.Println("GET - path", path, "fullpath:", fp)

	// open the file
	f, err := os.Open(fp)
	if err != nil {
		return nil, xerrors.Errorf("Unable to open file: %w", err)
	}
	return f, nil
}

// returns the absolute path
func (l *Local) fullpath(path string) string {
	return filepath.Join(l.basePath, path)
}
