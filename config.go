package main

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	CrawlingIntervalSeconds int `json:"crawling_interval_seconds"`
}

func LoadConfig(filename string) (*Configuration, error) {
	configFile, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Error while reading the config file: %s", err)
		return nil, err
	}

	var config Configuration
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Printf("Error while parsing the config file: %s", err)
		return nil, err
	}

	return &config, nil
}
