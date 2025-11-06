// Package config contains the configuration for the application.
package config

import (
	"github.com/BurntSushi/toml"
)

// Config defines the structure of config.toml
type Config struct {
	App      AppConfig      `toml:"app"`
	Server   ServerConfig   `toml:"server"`
	Logging  LoggingConfig  `toml:"logging"`
	Domains  DomainsConfig  `toml:"domains"`
	Database DatabaseConfig `toml:"database"`
}

type AppConfig struct {
	Name    string `toml:"name"`
	Version string `toml:"version"`
	Debug   bool   `toml:"debug"`
}

type ServerConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type LoggingConfig struct {
	Level string `toml:"level"`
	File  string `toml:"file"`
}

type DomainsConfig struct {
	Aliases []string `toml:"aliases"`
}

type DatabaseConfig struct {
	Path string `toml:"path"`
}

// LoadConfig loads the configuration from a TOML file path
func LoadConfig(path string) (*Config, error) {
	var conf Config
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
