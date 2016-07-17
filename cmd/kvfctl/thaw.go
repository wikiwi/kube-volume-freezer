/*
 * Copyright (C) 2016 wikiwi.io
 *
 * This software may be modified and distributed under the terms
 * of the MIT license. See the LICENSE file for details.
 */

package main

import (
	"fmt"

	"github.com/wikiwi/kube-volume-freezer/pkg/client"
	"github.com/wikiwi/kube-volume-freezer/pkg/util/validation"
)

type thawCommand struct {
}

func (cmd *thawCommand) Execute(args []string) error {
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
	_, err = client.Volumes().Thaw(globalOptions.Namespace, podName, volumeName)
	return err
}

func init() {
	parser.AddCommand("thaw",
		"Thaw Pod Volume",
		"Thaw Pod Volume",
		new(thawCommand))
}
