package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	APIURL     string `json:"api_url"`
	TunnelURL  string `json:"tunnel_url"`
	EncryptKey string `json:"encrypt_key"`
}

func configFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".zengate-config.json"), nil
}

func LoadConfig() (*Config, error) {
	path, err := configFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Config{}, nil
		}
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	return &cfg, err
}

func SaveConfig(cfg *Config) error {
	path, err := configFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

// ValidateEncryptKey checks if encryptKey is base64 and 32 bytes decoded
func ValidateEncryptKey(key string) error {
	key = strings.TrimSpace(key)
	decoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return errors.New("encrypt_key must be valid base64")
	}
	if len(decoded) != 32 {
		return errors.New("encrypt_key must decode to exactly 32 bytes")
	}
	return nil
}
