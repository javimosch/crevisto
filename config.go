package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	defaultAPIBase = "https://crevisto.intrane.fr"
	appName        = "crevisto"
)

// Config is the persistent CLI config stored at ~/.crevisto/config.json
type Config struct {
	APIToken string `json:"api_token"`
	APIBase  string `json:"api_base"`
}

func configDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".crevisto")
}

func configPath() string {
	return filepath.Join(configDir(), "config.json")
}

// loadConfig reads the config file, returning defaults if not found.
func loadConfig() *Config {
	c := &Config{APIBase: defaultAPIBase}
	data, err := os.ReadFile(configPath())
	if err == nil {
		json.Unmarshal(data, c)
	}
	if c.APIBase == "" {
		c.APIBase = defaultAPIBase
	}
	return c
}

// saveConfig writes the config file, creating the directory if needed.
func saveConfig(c *Config) error {
	dir := configDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	data, _ := json.MarshalIndent(c, "", "  ")
	return os.WriteFile(configPath(), data, 0600)
}

// getToken returns the stored bearer token, or empty string.
func getToken() string {
	return loadConfig().APIToken
}

// setToken saves the bearer token to config.
func setToken(token string) error {
	c := loadConfig()
	c.APIToken = token
	return saveConfig(c)
}
