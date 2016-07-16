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

// Options is the structure of the config file
type Options struct {
	ServerBackupFolder string `yaml:"ServerBackupFolder"`
	ServerPort         string `yaml:"ServerPort"`
}

type dependencies struct{}

func LoadConf() {
	b, err := ioutil.ReadFile(ConfFile)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(b, &Conf)
	if err != nil {
		log.Fatal(err)
	}
}
