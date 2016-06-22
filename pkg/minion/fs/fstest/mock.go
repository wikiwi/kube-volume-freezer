package fstest

import (
	"os"

	"github.com/stretchr/testify/mock"

	"github.com/wikiwi/kube-volume-freezer/pkg/minion/fs"
)

var _ fs.FileSystem = new(MockFS)

// MockFS is a filesystem mock.
type MockFS struct {
	mock.Mock
}

// ReadDir is mocked function.
func (fs *MockFS) ReadDir(path string) ([]os.FileInfo, error) {
	args := fs.Called(path)
	return args.Get(0).([]os.FileInfo), args.Error(1)
}

// DirExists is mocked function.
func (fs *MockFS) DirExists(path string) (bool, error) {
	args := fs.Called(path)
	return args.Bool(0), args.Error(1)
}

// Freeze is mocked function.
func (fs *MockFS) Freeze(path string) error {
	args := fs.Called(path)
	return args.Error(1)
}

// Thaw is mocked function.
func (fs *MockFS) Thaw(path string) error {
	args := fs.Called(path)
	return args.Error(1)
}
