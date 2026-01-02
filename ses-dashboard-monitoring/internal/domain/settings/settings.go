package settings

import (
	"context"
	"time"
)

type Setting struct {
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Description string    `json:"description"`
	IsEncrypted bool      `json:"is_encrypted"`
	UpdatedBy   int       `json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AWSConfig struct {
	Enabled   bool   `json:"enabled"`
	Region    string `json:"region"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

type TimezoneConfig struct {
	Timezone string `json:"timezone"`
}

type Repository interface {
	Get(ctx context.Context, key string) (*Setting, error)
	Set(ctx context.Context, key, value string, updatedBy int) error
	GetAll(ctx context.Context) ([]*Setting, error)
	GetAWSConfig(ctx context.Context) (*AWSConfig, error)
	TestAWSConnection(ctx context.Context, config *AWSConfig) error
	GetTimezoneConfig(ctx context.Context) (*TimezoneConfig, error)
}