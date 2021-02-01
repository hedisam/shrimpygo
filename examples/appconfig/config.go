package appconfig

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	APIKey    string `json:"api_key"`
	SecretKey string `json:"secret_key"`
}

func Read(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read app config file: %w", err)
	}

	var config *Config
	err = json.NewDecoder(f).Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode the app config json: %w", err)
	}

	return config, nil
}
