package fs

import (
	"os"
	"os/exec"

	"github.com/spf13/afero"

	"github.com/wikiwi/kube-volume-freezer/pkg/log"
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
	cmd := exec.Command("/sbin/fsfreeze", "-f", path)
	out, err := cmd.CombinedOutput()
	log := log.Instance().WithField("path", path).WithField("action", "freeze")
	if len(out) > 0 {
		log.Debugf("fsfreeze: %s", out)
	} else {
		log.Debugf("fsfreeze", out)
	}
	return err
}

func (fs *osFS) Thaw(path string) error {
	cmd := exec.Command("/sbin/fsfreeze", "-u", path)
	out, err := cmd.CombinedOutput()
	log := log.Instance().WithField("path", path).WithField("action", "thaw")
	if len(out) > 0 {
		log.Debugf("fsfreeze: %s", out)
	} else {
		log.Debugf("fsfreeze", out)
	}
	return err
}

// Default returns default OS File System.
func New() FileSystem {
	return &osFS{afero.Afero{Fs: afero.NewOsFs()}}
}
