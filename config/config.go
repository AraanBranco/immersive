package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Prefix       string   `json:"prefix"`
	DiscordToken string   `json:"discord_token"`
	OwnerId      string   `json:"owner_id"`
	UseSharding  bool     `json:"user_sharding"`
	ShardID      int      `json:"shard_id"`
	ShardCount   int      `json:"shard_count"`
	Cities       []string `json:"cities"`
	LocalDB      string   `json:"localdb"`
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
