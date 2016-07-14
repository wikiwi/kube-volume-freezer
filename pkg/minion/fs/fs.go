package fs

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/Sirupsen/logrus"
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
	log := log.Instance().WithField("path", path).WithField("action", "freeze")
	err := runExec(log, "/bin/sync")
	if err != nil {
		return err
	}
	return runExec(log, "/sbin/fsfreeze", "-f", path)
}

func (fs *osFS) Thaw(path string) error {
	log := log.Instance().WithField("path", path).WithField("action", "thaw")
	return runExec(log, "/sbin/fsfreeze", "-u", path)
}

// New returns default OS File System.
func New() FileSystem {
	return &osFS{afero.Afero{Fs: afero.NewOsFs()}}
}

func runExec(log *logrus.Entry, cmd string, flags ...string) error {
	c := exec.Command(cmd, flags...)
	var combinedOutput bytes.Buffer
	var errOut bytes.Buffer
	c.Stdout = &combinedOutput
	c.Stderr = io.MultiWriter(&combinedOutput, &errOut)
	err := c.Run()
	if combinedOutput.Len() > 0 {
		log.Debugf("%s: %s", cmd, combinedOutput.String())
	} else {
		log.Debugf("%s", cmd)
	}
	if err != nil {
		return errors.New(strings.TrimSpace(errOut.String()))
	}
	return nil
}
