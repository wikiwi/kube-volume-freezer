package fstest

import (
	"github.com/spf13/afero"

	"github.com/wikiwi/kube-volume-freezer/pkg/log"
)

// NewFake returns a new In-Memory-Filesystem with predefined directories.
func NewFake(dirs []string) *FakeFS {
	a := &FakeFS{Afero: afero.Afero{Fs: afero.NewMemMapFs()}}
	for _, d := range dirs {
		err := a.MkdirAll(d, 0755)
		if err != nil {
			log.Instance().Panic(err)
		}
	}
	return a
}

type FakeFS struct {
	afero.Afero
	Frozen []string
	Thawed []string
}

func (fs *FakeFS) Freeze(path string) error {
	fs.Frozen = append(fs.Frozen, path)
	return nil
}

func (fs *FakeFS) Thaw(path string) error {
	fs.Thawed = append(fs.Thawed, path)
	return nil
}
