package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFile = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (C *Config) SetUser(name string) error {
	C.CurrentUserName = name
	return Write(*C)
}

func Read() (Config, error) {
	path, err := getFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("ERROR getting filepath: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("ERROR reading: %w", err)
	}

	var structure Config
	if err := json.Unmarshal(data, &structure); err != nil {
		return Config{}, fmt.Errorf("ERROR unmarshalling: %w", err)
	}
	return structure, nil
}

func Write(cfg Config) error {
	fullPath, err := getFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}
	return nil
}

func getFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fullPath := filepath.Join(homeDir, configFile)
	return fullPath, nil
}
