/*
 * Copyright (C) 2016 wikiwi.io
 *
 * This software may be modified and distributed under the terms
 * of the MIT license. See the LICENSE file for details.
 */

package main

import (
	"fmt"

	"github.com/wikiwi/kube-volume-freezer/pkg/version"
)

type versionCommand struct {
}

func (x *versionCommand) Execute(args []string) error {
	fmt.Println(version.Version)
	return nil
}

func init() {
	parser.AddCommand("version",
		"Show version",
		"Print version to console.",
		new(versionCommand))
}
