package fs

import (
	"os"

	"github.com/spf13/afero"
)

// FileSystem is a abstraction of the file system with minimal required
// functionality.
type FileSystem interface {
	ReadDir(path string) ([]os.FileInfo, error)
	DirExists(path string) (bool, error)
	Freeze(path string) error
	Thaw(path string) error
}

type osFS struct {
	afero.Afero
}

func (fs *osFS) Freeze(path string) error {
	return nil
}

func (fs *osFS) Thaw(path string) error {
	return nil
}

// Default returns default OS File System.
func New() FileSystem {
	return &osFS{afero.Afero{Fs: afero.NewOsFs()}}
}
