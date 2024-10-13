package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the configuration for the application.
// It contains fields for storing the GROQ API key and commit prompt.
type Config struct {
	GROQAPIKey   string `json:"GROQ_APIKEY"`
	CommitPrompt string `json:"COMMIT_PROMPT"`
}

const (
	// configFileName is the name of the configuration file.
	configFileName = ".ai-commit"
)

var (
	// configFilePath holds the full path to the configuration file.
	configFilePath string
	// ErrUnknownKey is the error returned when an unknown configuration key is requested.
	ErrUnknownKey = fmt.Errorf("unknown config key")
)

func init() {
	// Initializes the configuration file path based on the user's home directory.
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("Failed to get user home directory: %v", err))
	}
	configFilePath = filepath.Join(homeDir, configFileName)
}

// loadConfig loads the configuration from the file.
// If the file does not exist, it returns an empty Config struct.
func loadConfig() (Config, error) {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, nil
		}
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	// Unmarshals the JSON data from the file into a Config struct.
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("failed to parse config file: %w", err)
	}
	return config, nil
}

// saveConfig saves the given Config struct to the configuration file.
// It writes the file with permission 0600 to ensure only the user can read/write it.
func saveConfig(config Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Writes the configuration data to the config file.
	if err := os.WriteFile(configFilePath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

// SetConfig updates the configuration for the given key with the specified value.
// It first loads the existing configuration, modifies the key, and then saves the updated config.
func SetConfig(key, value string) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	// Updates the corresponding field based on the provided key.
	switch key {
	case "GROQ_APIKEY":
		config.GROQAPIKey = value
	case "COMMIT_PROMPT":
		config.CommitPrompt = value
	default:
		// Returns an error if the key is not recognized.
		return fmt.Errorf("%w: %s", ErrUnknownKey, key)
	}

	// Saves the updated configuration.
	if err := saveConfig(config); err != nil {
		return err
	}

	fmt.Printf("Configuration updated: %s=%s\n", key, value)
	return nil
}

// GetConfig retrieves the value of the specified configuration key.
// It loads the current configuration and returns the value corresponding to the given key.
func GetConfig(key string) (string, error) {
	config, err := loadConfig()
	if err != nil {
		return "", err
	}

	// Returns the value based on the key or an error if the key is unknown.
	switch key {
	case "GROQ_APIKEY":
		return config.GROQAPIKey, nil
	case "COMMIT_PROMPT":
		return config.CommitPrompt, nil
	default:
		return "", fmt.Errorf("%w: %s", ErrUnknownKey, key)
	}
}

// GetConfigPath returns the full path to the configuration file.
func GetConfigPath() string {
	return configFilePath
}
