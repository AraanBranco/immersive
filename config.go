package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

)

type Config struct {
	Prefix        string `json:"prefix"`
	DiscordToken  string `json:"discord_token"`
	OwnerId       string `json:"owner_id"`
}

func LoadConfig(filename string) *Config {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("error loading config,", err)
		return nil
	}
	var conf Config
	json.Unmarshal(body, &conf)
	return &conf
}