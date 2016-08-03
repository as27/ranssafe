package main

import (
	"flag"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// ConfFile is the path to the configuration file
var ConfFile = flag.String("conf", "conf.yaml", "Path to the conf file")

// Conf stores the configuration
var Conf Options

// Pack defines a package on the client side
type Pack struct {
	Name string `yaml:"Name"`
	Path string `yaml:"Path"`
}

// Options is the structure of the config file
type Options struct {
	ServerURL string   `yaml:"ServerURL"`
	Packages  []Pack   `yaml:"Packages"`
	SkipDir   []string `yaml:"SkipDir"`
	MaxSize   int64    `yaml:"MaxSize"`
}

func init() {
	flag.Parse()
	b, err := ioutil.ReadFile(*ConfFile)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(b, &Conf)
	if err != nil {
		log.Fatal(err)
	}
}
