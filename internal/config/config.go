package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the configuration for the application,
// which includes the GROQ API key and the commit prompt.
type Config struct {
	GROQ_APIKEY   string `json:"GROQ_APIKEY"`   // API key for GROQ services.
	COMMIT_PROMPT string `json:"COMMIT_PROMPT"` // Commit prompt for generating commit messages.
}

var configFilePath string

// init initializes the config file path to store the configuration
// in the user's home directory with the filename ".ai-commit".
func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	configFilePath = filepath.Join(homeDir, ".ai-commit")
}

// loadConfig loads the configuration from the file.
// It returns an empty Config if the file doesn't exist.
func loadConfig() Config {
	var config Config

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		// If the config file does not exist, return an empty Config struct.
		if os.IsNotExist(err) {
			return Config{}
		}
		panic(err)
	}

	// Unmarshal the JSON data into the Config struct.
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	return config
}

// saveConfig saves the given Config struct to the configuration file in JSON format.
func saveConfig(config Config) {
	// Marshal the Config struct into a pretty-printed JSON format.
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		panic(err)
	}

	// Write the JSON data to the config file with 644 permissions.
	err = os.WriteFile(configFilePath, data, 0o644)
	if err != nil {
		panic(err)
	}
}

// SetConfig updates the configuration for the given key with the specified value
// and saves the updated configuration back to the file.
func SetConfig(key, value string) {
	config := loadConfig()

	// Update the appropriate configuration field based on the key.
	switch key {
	case "GROQ_APIKEY":
		config.GROQ_APIKEY = value
	case "COMMIT_PROMPT":
		config.COMMIT_PROMPT = value
	default:
		fmt.Printf("Unknown config key: %s\n", key)
		return
	}

	// Save the updated configuration to the file.
	saveConfig(config)
	fmt.Printf("Configuration updated: %s=%s\n", key, value)
}

// GetConfig retrieves the value of the specified configuration key from the config file.
// It returns "Not Set" if the key is unknown.
func GetConfig(key string) string {
	config := loadConfig()

	// Return the corresponding value based on the key.
	switch key {
	case "GROQ_APIKEY":
		return config.GROQ_APIKEY
	case "COMMIT_PROMPT":
		return config.COMMIT_PROMPT
	default:
		return "Not Set"
	}
}

// GetConfigPath returns the path to the configuration file.
func GetConfigPath() string {
	return configFilePath
}
