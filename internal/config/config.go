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

// Define a nice "Stringer" method to allow the Config construct to be printed as a string if called from functions like fmt.PrintLn()
func (c Config) String() string {
	return fmt.Sprintf("Database Connection URL: '%s' | User: '%s'", c.DbURL, c.CurrentUserName)
}

// Method to set username to a supplied one in the config file
func (c Config) SetUser(name string) error {
	newConfig, err := Read()
	if err != nil {
		return err
	}
	newConfig.CurrentUserName = name
	write(newConfig)
	return nil
}

// Helper function to obtain the path to the config file
func getConfigFilePath() (string, error) {
	configFilePath, err := os.UserHomeDir()
	if err != nil {
		return "Config file not found", err
	}
	return filepath.Join(configFilePath, configFileName), nil

}

// Read the config json file from the given path and return a Config struct
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

// Write the given Config strut out to the config json on disk
func write(cfg Config) error {
	configFilePath, _ := getConfigFilePath()
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	os.WriteFile(configFilePath, jsonData, 0644)
	return nil
}
