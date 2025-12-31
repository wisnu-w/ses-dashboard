package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App struct {
		Name          string `yaml:"name"`
		Env           string `yaml:"env"`
		Port          int    `yaml:"port"`
		LogBody       bool   `yaml:"log_body"`
		JWTSecret     string `yaml:"jwt_secret"`
		EnableSwagger bool   `yaml:"enable_swagger"`
	} `yaml:"app"`

	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`
}

func Load(path string) (*Config, error) {
	cfg := &Config{}

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(b, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
