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

func (c Config) String() string {
	return fmt.Sprintf("Database Connection URL: '%s' | User: '%s'", c.DbURL, c.CurrentUserName)
}

func Read() (Config, error) {
	configRtn := Config{}
	configFile, err := getConfigFilePath()
	// fmt.Println(configFile)
	if err != nil {
		return configRtn, err
	}
	rawConfig, err := os.ReadFile(configFile)
	if err != nil {
		return configRtn, err
	}
	json.Unmarshal(rawConfig, &configRtn)
	return configRtn, err
}

func SetUser(name string) error {
	newConfig, err := Read()
	if err != nil {
		return err
	}
	newConfig.CurrentUserName = name
	write(newConfig)
	return nil
}

func getConfigFilePath() (string, error) {
	configFilePath, err := os.UserHomeDir()
	// fmt.Println(configFilePath)
	// fmt.Println(filepath.Join(configFilePath, configFileName))
	if err != nil {
		return "Config file not found", err
	}
	return filepath.Join(configFilePath, configFileName), nil

}

func write(cfg Config) error {
	configFilePath, _ := getConfigFilePath()
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	os.WriteFile(configFilePath, jsonData, 0644)
	return nil
}
