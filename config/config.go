package config

import (
	"errors"
	"georgslauf/auth"
	"georgslauf/persistence"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigData struct {
	Server struct {
		Host string `yaml:"host" binding:"required"`
		Port int64  `yaml:"port" binding:"required"`
	} `yaml:"server"`
	UploadDir string                     `yaml:"uploadDir"`
	Database  persistence.DatabaseConfig `yaml:"database"`
	OAuth     auth.OAuthConfig           `yaml:"oauth"`
	Security  SecurityConfig             `yaml:"security"`
}

type SecurityConfig struct {
	CSRFAuthKey []byte `yaml:"csrfAuthKey"`
	SessionKey  []byte `yaml:"sessionKey"`
}

var (
	ErrorValidation = errors.New("config validation failed")
)

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

	if len(cfg.Security.CSRFAuthKey) != 32 {
		return nil, ErrorValidation
	}

	return &cfg, nil
}
