package repository

import (
	"context"
	"database/sql"
	"fmt"

	"ses-monitoring/internal/domain/settings"
)

type settingsRepo struct {
	db *sql.DB
}

func NewSettingsRepository(db *sql.DB) settings.Repository {
	return &settingsRepo{db: db}
}

func (r *settingsRepo) Get(ctx context.Context, key string) (*settings.Setting, error) {
	query := `SELECT key, value, COALESCE(description, '') as description, is_encrypted, updated_by, created_at, updated_at FROM app_settings WHERE key = $1`

	s := &settings.Setting{}
	err := r.db.QueryRowContext(ctx, query, key).Scan(
		&s.Key, &s.Value, &s.Description, &s.IsEncrypted, &s.UpdatedBy, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *settingsRepo) Set(ctx context.Context, key, value string, updatedBy int) error {
	if updatedBy == 0 {
		return fmt.Errorf("user ID is required for settings update")
	}

	descriptions := map[string]string{
		"aws_enabled":       "Enable/disable AWS SES integration",
		"aws_region":        "AWS region for SES service",
		"aws_access_key":    "AWS access key for SES authentication",
		"aws_secret_key":    "AWS secret key for SES authentication",
		"retention_days":    "Number of days to retain event logs (0 = never delete)",
		"retention_enabled": "Enable/disable automatic log retention cleanup",
	}

	description := descriptions[key]
	if description == "" {
		description = "Application setting"
	}

	query := `
		INSERT INTO app_settings (key, value, description, updated_by, updated_at) 
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (key) 
		DO UPDATE SET value = $2, updated_by = $4, updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query, key, value, description, updatedBy)
	return err
}

func (r *settingsRepo) GetAll(ctx context.Context) ([]*settings.Setting, error) {
	query := `SELECT key, value, COALESCE(description, '') as description, is_encrypted, updated_by, created_at, updated_at FROM app_settings ORDER BY key`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settingsList []*settings.Setting
	for rows.Next() {
		s := &settings.Setting{}
		err := rows.Scan(&s.Key, &s.Value, &s.Description, &s.IsEncrypted, &s.UpdatedBy, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		settingsList = append(settingsList, s)
	}
	return settingsList, nil
}

func (r *settingsRepo) GetAWSConfig(ctx context.Context) (*settings.AWSConfig, error) {
	config := &settings.AWSConfig{}

	if enabled, err := r.Get(ctx, "aws_enabled"); err == nil {
		config.Enabled = enabled.Value == "true"
	}

	if region, err := r.Get(ctx, "aws_region"); err == nil {
		config.Region = region.Value
	}

	if accessKey, err := r.Get(ctx, "aws_access_key"); err == nil {
		config.AccessKey = accessKey.Value
	}

	if secretKey, err := r.Get(ctx, "aws_secret_key"); err == nil {
		config.SecretKey = secretKey.Value
	}

	return config, nil
}

func (r *settingsRepo) TestAWSConnection(ctx context.Context, config *settings.AWSConfig) error {
	// This will be implemented with actual AWS SDK
	if config.AccessKey == "" || config.SecretKey == "" {
		return fmt.Errorf("AWS credentials are required")
	}
	return nil
}
