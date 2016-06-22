package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"

	"github.com/wikiwi/kube-volume-freezer/pkg/log"
	"github.com/wikiwi/kube-volume-freezer/pkg/minion"
	"github.com/wikiwi/kube-volume-freezer/pkg/version"
)

var opts struct {
	Listen  string `long:"listen" default:"0.0.0.0:8080" env:"KVF_LISTEN" description:"Address to listen on"`
	Token   string `short:"t" long:"token" env:"KVF_TOKEN" description:"Use given token for api user authentication"`
	Verbose bool   `short:"v" long:"verbose" description:"Turn on verbose logging"`
	Version bool   `long:"version" description:"Show version"`
}

func main() {
	var parser = flags.NewParser(&opts, flags.Default)
	parser.Name = "kvf-minion"
	parser.LongDescription = "Run minion server. This should be run in a Kubernetes Pod on each Node of the Cluster."

	_, err := parser.Parse()
	if err != nil {
		if e2, ok := err.(*flags.Error); ok && e2.Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}
	log.SetupAndHarmonize(opts.Verbose)
	if opts.Version {
		fmt.Println(version.Version)
		os.Exit(0)
	}
	server, err := minion.NewRestServer(&minion.Options{
		Token: opts.Token,
	})
	if err != nil {
		log.Instance().Panic(err)
	}
	log.Instance().Panic(server.ListenAndServe(opts.Listen))
}
