package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the configuration for the application.
type Config struct {
	GROQAPIKey   string `json:"GROQ_APIKEY"`
	CommitPrompt string `json:"COMMIT_PROMPT"`
}

const (
	configFileName = ".ai-commit"
)

var (
	configFilePath string
	ErrUnknownKey  = fmt.Errorf("unknown config key")
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("Failed to get user home directory: %v", err))
	}
	configFilePath = filepath.Join(homeDir, configFileName)
}

// loadConfig loads the configuration from the file.
func loadConfig() (Config, error) {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, nil
		}
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("failed to parse config file: %w", err)
	}
	return config, nil
}

// saveConfig saves the given Config struct to the configuration file.
func saveConfig(config Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configFilePath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

// SetConfig updates the configuration for the given key with the specified value.
func SetConfig(key, value string) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	switch key {
	case "GROQ_APIKEY":
		config.GROQAPIKey = value
	case "COMMIT_PROMPT":
		config.CommitPrompt = value
	default:
		return fmt.Errorf("%w: %s", ErrUnknownKey, key)
	}

	if err := saveConfig(config); err != nil {
		return err
	}

	fmt.Printf("Configuration updated: %s=%s\n", key, value)
	return nil
}

// GetConfig retrieves the value of the specified configuration key.
func GetConfig(key string) (string, error) {
	config, err := loadConfig()
	if err != nil {
		return "", err
	}

	switch key {
	case "GROQ_APIKEY":
		return config.GROQAPIKey, nil
	case "COMMIT_PROMPT":
		return config.CommitPrompt, nil
	default:
		return "", fmt.Errorf("%w: %s", ErrUnknownKey, key)
	}
}

// GetConfigPath returns the path to the configuration file.
func GetConfigPath() string {
	return configFilePath
}

