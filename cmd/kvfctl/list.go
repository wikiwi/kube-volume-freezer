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

type listCommand struct {
}

func (cmd *listCommand) Execute(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Error: Please specify Pod and Volume.")
	}
	if len(args) > 1 {
		return fmt.Errorf("Error: Unexpected argument %s", args[1:])
	}

	podName := args[0]

	if issues := validation.ValidateQualitfiedName(podName); len(issues) > 0 {
		return fmt.Errorf("Error: Invalid Pod Name %s", issues)
	}

	options := &client.Options{Token: globalOptions.Token}
	client, err := client.New(globalOptions.Address, options)
	if err != nil {
		return err
	}
	ls, err := client.Volumes().List(globalOptions.Namespace, podName)
	for _, item := range ls.Items {
		fmt.Println(item)
	}
	return err
}

func init() {
	parser.AddCommand("list",
		"List Pod Volumes",
		"List Pod Volumes",
		new(listCommand))
}
