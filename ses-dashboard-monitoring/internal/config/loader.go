package config

import (
	"os"
	"strconv"

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

	AWS struct {
		Region    string `yaml:"region"`
		AccessKey string `yaml:"access_key"`
		SecretKey string `yaml:"secret_key"`
	} `yaml:"aws"`
}

func Load(path string) (*Config, error) {
	cfg := &Config{}

	// First, try to load from environment variables (Docker priority)
	cfg.App.Name = getEnv("APP_NAME", "")
	cfg.App.Env = getEnv("APP_ENV", "")
	cfg.App.Port = getEnvInt("APP_PORT", 0)
	cfg.App.JWTSecret = getEnv("JWT_SECRET", "")
	cfg.App.EnableSwagger = getEnvBool("ENABLE_SWAGGER", false)

	cfg.Database.Host = getEnv("DB_HOST", "")
	cfg.Database.Port = getEnvInt("DB_PORT", 0)
	cfg.Database.Name = getEnv("DB_NAME", "")
	cfg.Database.User = getEnv("DB_USER", "")
	cfg.Database.Password = getEnv("DB_PASSWORD", "")
	cfg.Database.SSLMode = getEnv("DB_SSLMODE", "")

	cfg.AWS.Region = getEnv("AWS_REGION", "")
	cfg.AWS.AccessKey = getEnv("AWS_ACCESS_KEY", "")
	cfg.AWS.SecretKey = getEnv("AWS_SECRET_KEY", "")

	// If environment variables are not set, fallback to YAML file
	if cfg.App.Name == "" || cfg.Database.Host == "" {
		if b, err := os.ReadFile(path); err == nil {
			yamlCfg := &Config{}
			if err := yaml.Unmarshal(b, yamlCfg); err == nil {
				// Use YAML values only if env vars are empty
				if cfg.App.Name == "" { cfg.App.Name = yamlCfg.App.Name }
				if cfg.App.Env == "" { cfg.App.Env = yamlCfg.App.Env }
				if cfg.App.Port == 0 { cfg.App.Port = yamlCfg.App.Port }
				if cfg.App.JWTSecret == "" { cfg.App.JWTSecret = yamlCfg.App.JWTSecret }
				if os.Getenv("ENABLE_SWAGGER") == "" { cfg.App.EnableSwagger = yamlCfg.App.EnableSwagger }
				cfg.App.LogBody = yamlCfg.App.LogBody

				if cfg.Database.Host == "" { cfg.Database.Host = yamlCfg.Database.Host }
				if cfg.Database.Port == 0 { cfg.Database.Port = yamlCfg.Database.Port }
				if cfg.Database.Name == "" { cfg.Database.Name = yamlCfg.Database.Name }
				if cfg.Database.User == "" { cfg.Database.User = yamlCfg.Database.User }
				if cfg.Database.Password == "" { cfg.Database.Password = yamlCfg.Database.Password }
				if cfg.Database.SSLMode == "" { cfg.Database.SSLMode = yamlCfg.Database.SSLMode }

				if cfg.AWS.Region == "" { cfg.AWS.Region = yamlCfg.AWS.Region }
				if cfg.AWS.AccessKey == "" { cfg.AWS.AccessKey = yamlCfg.AWS.AccessKey }
				if cfg.AWS.SecretKey == "" { cfg.AWS.SecretKey = yamlCfg.AWS.SecretKey }
			}
		}
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}