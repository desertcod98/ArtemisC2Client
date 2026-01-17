package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	AgentId        string `json:"agent_id"`
	BeaconInterval int    `json:"beacon_interval"` // in seconds
}

const filepath = "cfg"

func SaveConfig(cfg *Config) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer f.Close()
	return json.NewEncoder(f).Encode(cfg)
}

func LoadConfig() (*Config, error) {
	var cfg Config
	f, err := os.Open(filepath)
	if err != nil {
		return &cfg, err
	}

	defer f.Close()
	err = json.NewDecoder(f).Decode(&cfg)
	return &cfg, err
}
