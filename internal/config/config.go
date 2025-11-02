package config

import (
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

const (
	ConfigFileName = ".smart-clipboard.yaml"
)

type Config struct {
	ClipboardHistorySize int  `yaml:"clipboard_history_size"`
	DebugMode            bool `yaml:"debug_mode"`
	RunAtStartup         bool `yaml:"run_at_startup"`
}

func DefaultConfig() *Config {
	return &Config{
		ClipboardHistorySize: 50,
		DebugMode:            false,
		RunAtStartup:         false,
	}
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		instance = loadConfig()
	})
	return instance
}

func loadConfig() *Config {
	config := DefaultConfig()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return config
	}

	configPath := filepath.Join(homeDir, ConfigFileName)

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			saveConfig(config, configPath)
			return config
		}
		return config
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return config
	}

	return config
}

func saveConfig(config *Config, configPath string) error {
	dir := filepath.Dir(configPath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, data, 0644)
	return err
}

func SaveConfig() error {
	config := GetConfig()
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, ConfigFileName)
	return saveConfig(config, configPath)
}

func SetClipboardHistorySize(size int) {
	GetConfig().ClipboardHistorySize = size
}

func SetDebugMode(debug bool) {
	GetConfig().DebugMode = debug
}

func SetRunAtStartup(runAtStartup bool) error {
	return setRunAtStartup(runAtStartup)
}

func setRunAtStartup(runAtStartup bool) error {
	config := GetConfig()
	config.RunAtStartup = runAtStartup

	// Update system startup settings
	err := updateStartupSettings(runAtStartup)
	if err != nil {
		return err
	}

	return SaveConfig()
}

func GetRunAtStartup() bool {
	return GetConfig().RunAtStartup
}
