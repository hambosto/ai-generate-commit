package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	GROQ_APIKEY   string `json:"GROQ_APIKEY"`
	COMMIT_PROMPT string `json:"COMMIT_PROMPT"`
}

var configFilePath string

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	configFilePath = filepath.Join(homeDir, ".ai-commit")
}

func loadConfig() Config {
	var config Config

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}
		}
		panic(err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	return config
}

func saveConfig(config Config) {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(configFilePath, data, 0o644)
	if err != nil {
		panic(err)
	}
}

func SetConfig(key, value string) {
	config := loadConfig()

	switch key {
	case "GROQ_APIKEY":
		config.GROQ_APIKEY = value
	case "COMMIT_PROMPT":
		config.COMMIT_PROMPT = value
	default:
		fmt.Printf("Unknown config key: %s\n", key)
		return
	}

	saveConfig(config)
	fmt.Printf("Configuration  updated: %s=%s\n", key, value)
}

func GetConfig(key string) string {
	config := loadConfig()

	switch key {
	case "GROQ_APIKEY":
		return config.GROQ_APIKEY
	case "COMMIT_PROMPT":
		return config.COMMIT_PROMPT
	default:
		return "Not Set"
	}
}

func GetConfigPath() string {
	return configFilePath
}
