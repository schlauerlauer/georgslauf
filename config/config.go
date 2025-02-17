package config

import (
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigData struct {
	Server struct {
		Host string `yaml:"host" binding:"required"`
		Port int64  `yaml:"port" binding:"required"`
	} `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

func NewConfig(path string) (*ConfigData, error) {
	var cfg ConfigData
	if path == "" {
		path = "./config.yaml"
	}

	fileContent, err := os.ReadFile(path)
	if err != nil {
		slog.Error("error opening config file", path, err)
		return nil, err
	}

	if err = yaml.Unmarshal(fileContent, &cfg); err != nil {
		slog.Error("error parsing config", "err", err)
		os.Exit(1)
	}

	return &cfg, nil
}
