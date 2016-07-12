package fs

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"

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

func (fs *osFS) fsfreeze(path, action string, flags ...string) error {
	args := append(flags, path)
	cmd := exec.Command("/sbin/fsfreeze", args...)
	var combinedOutput bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &combinedOutput
	cmd.Stderr = io.MultiWriter(&combinedOutput, &errOut)
	err := cmd.Run()
	log := log.Instance().WithField("path", path).WithField("action", action)
	if combinedOutput.Len() > 0 {
		log.Debugf("fsfreeze: %s", combinedOutput.String())
	} else {
		log.Debug("fsfreeze")
	}
	if err != nil {
		return errors.New(strings.TrimSpace(errOut.String()))
	}

	return nil
}

func (fs *osFS) Freeze(path string) error {
	return fs.fsfreeze(path, "freeze", "-f")
}

func (fs *osFS) Thaw(path string) error {
	return fs.fsfreeze(path, "thaw", "-u")
}

// New returns default OS File System.
func New() FileSystem {
	return &osFS{afero.Afero{Fs: afero.NewOsFs()}}
}
