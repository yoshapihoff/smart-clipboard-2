package config

import (
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

const (
	// ConfigFileName is the name of the configuration file
	ConfigFileName = ".smart-clipboard.yaml"
)

// Config holds the application configuration
type Config struct {
	// Размер списка с историей буфера обмена
	ClipboardHistorySize int `yaml:"clipboard_history_size"`

	// Режим дебага
	DebugMode bool `yaml:"debug_mode"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		ClipboardHistorySize: 50,
		DebugMode:            false,
	}
}

var (
	instance *Config
	once     sync.Once
)

// GetConfig returns the singleton configuration instance
func GetConfig() *Config {
	once.Do(func() {
		instance = loadConfig()
	})
	return instance
}

// loadConfig loads configuration from YAML file or creates default
func loadConfig() *Config {
	config := DefaultConfig()

	// Try to load from config file in user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return config
	}

	configPath := filepath.Join(homeDir, ConfigFileName)

	data, err := os.ReadFile(configPath)
	if err != nil {
		// If file doesn't exist, save default config and return it
		if os.IsNotExist(err) {
			saveConfig(config, configPath)
			return config
		}
		return config
	}

	// Parse YAML file
	err = yaml.Unmarshal(data, config)
	if err != nil {
		// If parsing fails, return default config
		return config
	}

	return config
}

// saveConfig saves the configuration to YAML file
func saveConfig(config *Config, configPath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(configPath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	// Marshal config to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	// Write to file
	err = os.WriteFile(configPath, data, 0644)
	return err
}

// SaveConfig saves the current configuration to file
func SaveConfig() error {
	config := GetConfig()
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, ConfigFileName)
	return saveConfig(config, configPath)
}

// SetClipboardHistorySize sets the clipboard history size
func SetClipboardHistorySize(size int) {
	GetConfig().ClipboardHistorySize = size
}

// SetDebugMode sets the debug mode
func SetDebugMode(debug bool) {
	GetConfig().DebugMode = debug
}