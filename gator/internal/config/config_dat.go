package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (C *Config) SetUser(name string) error {
	C.CurrentUserName = name

	configData, err := json.Marshal(C)
	if err != nil {
		return fmt.Errorf("ERROR marshalling: %w", err)
	}

	err = os.WriteFile(".gatorconfig.json", configData, 0644)
	if err != nil {
		return fmt.Errorf("ERROR writing: %w", err)
	}
	return nil
}

func Read(path string) (Config, error) {
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
