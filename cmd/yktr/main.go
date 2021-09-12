package main

import (
	"io/ioutil"
	"log"

	"github.com/pelletier/go-toml"
	"github.com/winebarrel/yktr"
)

func init() {
	log.SetFlags(log.LstdFlags)
}

func main() {
	opts := parseArgs()
	data, err := ioutil.ReadFile(opts.Config)

	if err != nil {
		log.Fatal(err)
	}

	config := &yktr.Config{}
	err = toml.Unmarshal(data, config)

	if err != nil {
		log.Fatal(err)
	}

	err = config.Validate()

	if err != nil {
		log.Fatal(err)
	}

	svr, err := yktr.NewServer(config)

	if err != nil {
		log.Fatal(err)
	}

	err = svr.Run()

	if err != nil {
		log.Fatal(err)
	}
}
