package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type ShrimpyConfig struct {
	APIKey    string `json:"api_key"`
	SecretKey string `json:"secret_key"`
}

func Read(path string) (*ShrimpyConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config *ShrimpyConfig
	err = json.NewDecoder(f).Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode the json config data: %w", err)
	}

	return config, nil
}
