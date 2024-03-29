package config

import (
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigInterface interface{}

type ConfigService struct {
	Config *ConfigData
}

type ConfigData struct {
	Server struct {
		Host string `yaml:"host" binding:"required"`
		Port string `yaml:"port" binding:"required"`
	} `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Auth     AuthConfig     `yaml:"auth"`
}

type DatabaseConfig struct {
	Path     string `yaml:"path"`
	Timezone string `yaml:"timezone"`
}

type AuthConfig struct {
	KratosLocalURL  string `yaml:"kratosLocal"`
	KratosPublicURL string `yaml:"kratosPublic"`
}

var _ ConfigInterface = &ConfigData{}

func NewConfig(path string) (*ConfigService, error) {
	var cfg *ConfigData
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

	return &ConfigService{
		Config: cfg,
	}, nil
}
