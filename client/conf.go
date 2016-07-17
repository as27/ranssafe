package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// ConfFile is the path to the configuration file
var ConfFile = "conf.yaml"

// Conf stores the configuration
var Conf Options

// Pack defines a package on the client side
type Pack struct {
	Name string `yaml:"Name"`
	Path string `yaml:"Path"`
}

// Options is the structure of the config file
type Options struct {
	ServerURL string `yaml:"ServerURL"`
	Packages  []Pack `yaml:"Packages"`
}

func init() {
	b, err := ioutil.ReadFile(ConfFile)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(b, &Conf)
	if err != nil {
		log.Fatal(err)
	}
}
