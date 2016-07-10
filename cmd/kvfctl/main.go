package main

import (
	"os"

	"github.com/jessevdk/go-flags"

	"github.com/wikiwi/kube-volume-freezer/pkg/log"
)

var globalOptions struct {
	Verbose bool `short:"v" long:"verbose" description:"Turn on verbose logging"`
}

var parser = flags.NewParser(&globalOptions, flags.Default)

func main() {
	parser.CommandHandler = func(command flags.Commander, args []string) error {
		log.SetupAndHarmonize(globalOptions.Verbose)
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
