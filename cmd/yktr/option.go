package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
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

	flag.StringVar(&opts.Config, "config", "", "config file")
	ver := flag.Bool("version", false, "print version")
	flag.Parse()

	if *ver {
		printVersionAndEixt()
	}

	if opts.Config == "" {
		exePath, err := os.Executable()

		if err != nil {
			log.Fatal(err)
		}

		opts.Config = path.Join(filepath.Dir(exePath), DefaultConfig)
	}

	return opts
}

func printVersionAndEixt() {
	fmt.Fprintln(os.Stderr, version)
	os.Exit(0)
}
