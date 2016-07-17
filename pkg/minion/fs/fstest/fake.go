/*
 * Copyright (C) 2016 wikiwi.io
 *
 * This software may be modified and distributed under the terms
 * of the MIT license. See the LICENSE file for details.
 */

package fstest

import (
	"github.com/spf13/afero"

	"github.com/wikiwi/kube-volume-freezer/pkg/log"
)

// NewFake creates a new instance of FakeFS.
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

// FakeFS is an In-Memory-Filesystem with predefined directories.
type FakeFS struct {
	afero.Afero
	Frozen []string
	Thawed []string
}

// Freeze adds path to field Frozen.
func (fs *FakeFS) Freeze(path string) error {
	fs.Frozen = append(fs.Frozen, path)
	return nil
}

// Thaw adds path to field Thawed.
func (fs *FakeFS) Thaw(path string) error {
	fs.Thawed = append(fs.Thawed, path)
	return nil
}
