package main

import (
	"fmt"

	"github.com/wikiwi/kube-volume-freezer/pkg/client"
	"github.com/wikiwi/kube-volume-freezer/pkg/util/validation"
)

type freezeCommand struct {
	Address   string `long:"address" default:"localhost:8080" env:"KVF_ADDRESS" description:"Address of kvf-master"`
	Namespace string `long:"namespace" default:"default" env:"KVF_NAMESPACE" description:"Namespace of Pod"`
	Token     string `short:"t" long:"token" env:"KVF_TOKEN" description:"Use given token for api user authentication"`
}

func (cmd *freezeCommand) Execute(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Error: Please specify Pod and Volume.")
	}
	if len(args) > 2 {
		return fmt.Errorf("Error: Unexpected argument %s", args[2:])
	}

	podName, volumeName := args[0], args[1]

	if issues := validation.ValidateQualitfiedName(cmd.Namespace); len(issues) > 0 {
		return fmt.Errorf("Error: Invalid Namespace %s", issues)
	}
	if issues := validation.ValidateQualitfiedName(podName); len(issues) > 0 {
		return fmt.Errorf("Error: Invalid Pod Name %s", issues)
	}
	if issues := validation.ValidateQualitfiedName(volumeName); len(issues) > 0 {
		return fmt.Errorf("Error: Invalid Volume Name %s", issues)
	}

	client, err := client.New(cmd.Address, cmd.Token, nil)
	if err != nil {
		return err
	}
	_, err = client.Volumes().Freeze(cmd.Namespace, podName, volumeName)
	return err
}

func init() {
	parser.AddCommand("freeze",
		"Freeze Pod Volume",
		"Freeze Pod Volume",
		new(freezeCommand))
}
