package suppression

import (
	"context"
	"time"
)

type SuppressionType string

const (
	SuppressionTypeBounce    SuppressionType = "bounce"
	SuppressionTypeComplaint SuppressionType = "complaint"
	SuppressionTypeManual    SuppressionType = "manual"
)

type AWSStatus string

const (
	AWSStatusUnknown       AWSStatus = "unknown"
	AWSStatusSuppressed    AWSStatus = "suppressed"
	AWSStatusNotSuppressed AWSStatus = "not_suppressed"
)

type SuppressionEntry struct {
	ID              int             `json:"id"`
	Email           string          `json:"email"`
	SuppressionType SuppressionType `json:"suppression_type"`
	Reason          string          `json:"reason"`
	AWSStatus       AWSStatus       `json:"aws_status"`
	IsActive        bool            `json:"is_active"`
	AddedBy         int             `json:"added_by"`
	AddedByName     string          `json:"added_by_name,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	SyncedAt        *time.Time      `json:"synced_at,omitempty"`
}

type Repository interface {
	// Local operations
	Add(ctx context.Context, entry *SuppressionEntry) error
	Remove(ctx context.Context, email string) error
	GetAll(ctx context.Context, limit, offset int) ([]*SuppressionEntry, error)
	GetCount(ctx context.Context) (int, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*SuppressionEntry, error)
	GetSearchCount(ctx context.Context, query string) (int, error)
	IsSupressed(ctx context.Context, email string) (bool, error)
	
	// AWS sync operations
	UpdateAWSStatus(ctx context.Context, email string, status AWSStatus) error
	GetUnsyncedEntries(ctx context.Context) ([]*SuppressionEntry, error)
	MarkAsSynced(ctx context.Context, email string) error
}