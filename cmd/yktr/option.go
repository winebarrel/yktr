package main

import (
	"flag"
	"fmt"
	"os"
)

var version string

const (
	DefaultConfig = "yktr.toml"
)

type Options struct {
	Config string
}

func parseArgs() *Options {
	opts := &Options{}

	flag.StringVar(&opts.Config, "config", DefaultConfig, "config file")
	ver := flag.Bool("version", false, "print version")
	flag.Parse()

	if *ver {
		printVersionAndEixt()
	}

	return opts
}

func printVersionAndEixt() {
	fmt.Fprintln(os.Stderr, version)
	os.Exit(0)
}
