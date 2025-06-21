package config

import (
	"flag"
	"log"

	"gopkg.in/yaml.v2"
	"os"
)

var cfg App

func Setup() {
	path := flag.String("config", "./config/config.yaml", "the absolute path of config.yaml")
	flag.Parse()

	content, err := os.ReadFile(*path)
	if err != nil {
		log.Fatalf("cfg ReadFile error: %v", err)
	}

	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		log.Fatalf("cfg Unmarshal error: %v", err)
	}
}
