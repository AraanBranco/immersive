package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Prefix       string `json:"prefix"`
	DiscordToken string `json:"discord_token"`
	OwnerId      string `json:"owner_id"`
	UseSharding  bool   `json:"owner_id"`
	ShardID      int    `json:"shard_id"`
	ShardCount   int    `json:"shard_count"`
}

func LoadConfig(filename string) *Config {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("error loading config,", err)
		return nil
	}
	var config Config
	json.Unmarshal(body, &config)
	return &config
}
