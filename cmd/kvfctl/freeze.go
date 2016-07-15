package main

import (
	"fmt"

	"github.com/wikiwi/kube-volume-freezer/pkg/client"
	"github.com/wikiwi/kube-volume-freezer/pkg/util/validation"
)

type freezeCommand struct {
}

func (cmd *freezeCommand) Execute(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("Error: Please specify Pod and Volume.")
	}
	if len(args) > 2 {
		return fmt.Errorf("Error: Unexpected argument %s", args[2:])
	}

	podName, volumeName := args[0], args[1]

	if issues := validation.ValidateQualitfiedName(podName); len(issues) > 0 {
		return fmt.Errorf("Error: Invalid Pod Name %s", issues)
	}
	if issues := validation.ValidateQualitfiedName(volumeName); len(issues) > 0 {
		return fmt.Errorf("Error: Invalid Volume Name %s", issues)
	}

	options := &client.Options{Token: globalOptions.Token}
	client, err := client.New(globalOptions.Address, options)
	if err != nil {
		return err
	}
	_, err = client.Volumes().Freeze(globalOptions.Namespace, podName, volumeName)
	return err
}

func init() {
	parser.AddCommand("freeze",
		"Freeze Pod Volume",
		"Freeze Pod Volume",
		new(freezeCommand))
}
