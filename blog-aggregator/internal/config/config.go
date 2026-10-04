package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DB_URL   string `json:"db_url"`
	Username string `json:"current_user_name"`
}

const CONFIG_FILENAME string = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", home, CONFIG_FILENAME), nil
}

func Read() (Config, error) {
	var config Config
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return config, err
	}
	file_bytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(file_bytes, &config)
	return config, err
}

func (c *Config) SetUser(username string) error {
	newConfig := c
	newConfig.Username = username
	configBytes, err := json.Marshal(newConfig)
	if err != nil {
		return err
	}
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	err = os.WriteFile(configFilePath, configBytes, os.ModePerm)
	if err != nil {
		return err
	}
	c = newConfig
	return nil
}
