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
