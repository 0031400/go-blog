package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	ListenAddr    string
	AdminAccount  string
	AdminPassword string
}

var config Config

func loalConfig() {
	data, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
}
