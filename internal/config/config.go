package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	configRtn := Config{}
	configFile, err := getConfigFilePath()
	if err != nil {
		rawConfig, err := os.ReadFile(configFile)
		if err != nil {
			json.Unmarshal(rawConfig, &configRtn)
		}
	}
	return configRtn, err
}

func getConfigFilePath() (string, error) {
	configFilePath, err := os.UserHomeDir()
	fmt.Println(configFilePath)
	fmt.Println(filepath.Join(configFilePath, configFileName))
	if err != nil {
		return filepath.Join(configFilePath, configFileName), nil
	}
	return "Config file not found", err
}
