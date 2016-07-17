/*
 * Copyright (C) 2016 wikiwi.io
 *
 * This software may be modified and distributed under the terms
 * of the MIT license. See the LICENSE file for details.
 */

// kvfctl is a command line client to the Master API Server.
package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"

	"github.com/wikiwi/kube-volume-freezer/pkg/log"
	"github.com/wikiwi/kube-volume-freezer/pkg/util/validation"
)

var globalOptions struct {
	Address   string `long:"address" default:"http://localhost:8080" env:"KVF_ADDRESS" description:"Address of kvf-master"`
	Namespace string `long:"namespace" default:"default" env:"KVF_NAMESPACE" description:"Namespace of Pod"`
	Token     string `short:"t" long:"token" env:"KVF_TOKEN" description:"Use given token for api user authentication"`
	Verbose   bool   `short:"v" long:"verbose" description:"Turn on verbose logging"`
}

var parser = flags.NewParser(&globalOptions, flags.Default)

func main() {
	parser.CommandHandler = func(command flags.Commander, args []string) error {
		log.SetupAndHarmonize(globalOptions.Verbose)
		if issues := validation.ValidateQualitfiedName(globalOptions.Namespace); len(issues) > 0 {
			return fmt.Errorf("Error: Invalid Namespace %s", issues)
		}
		return command.Execute(args)
	}

	parser.Name = "kvfctl"
	_, err := parser.Parse()
	if err != nil {
		if e2, ok := err.(*flags.Error); ok && e2.Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}
}
