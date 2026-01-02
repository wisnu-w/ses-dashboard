package services

import (
	"context"
	"log"
	"sync"
	"time"

	"ses-monitoring/internal/domain/models"
	"ses-monitoring/internal/domain/settings"
	"ses-monitoring/internal/infrastructure/aws"
	"ses-monitoring/internal/infrastructure/database"
)

type SyncService struct {
	settingsRepo settings.Repository
	dbRepo       *database.SuppressionRepository

	mu             sync.RWMutex
	lastSync       time.Time
	syncInProgress bool
}

func NewSyncService(
	settingsRepo settings.Repository,
	dbRepo *database.SuppressionRepository,
) *SyncService {
	return &SyncService{
		settingsRepo: settingsRepo,
		dbRepo:       dbRepo,
	}
}

// StartBackgroundSync memulai sync otomatis setiap 5 menit
func (s *SyncService) StartBackgroundSync(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	// Sync pertama setelah startup
	go func() {
		time.Sleep(10 * time.Second)
		_ = s.SyncNow(context.Background())
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			go s.SyncNow(context.Background())
		}
	}
}

// SyncNow melakukan sync manual
func (s *SyncService) SyncNow(ctx context.Context) error {
	s.mu.Lock()
	if s.syncInProgress {
		s.mu.Unlock()
		log.Println("Sync already in progress, skipping...")
		return nil
	}
	s.syncInProgress = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.syncInProgress = false
		s.lastSync = time.Now()
		s.mu.Unlock()
		log.Println("Sync process completed")
	}()

	// 🔥 AMBIL CONFIG TERBARU DARI DB
	cfg, err := s.settingsRepo.GetAWSConfig(ctx)
	if err != nil {
		log.Printf("Failed to load AWS config: %v", err)
		return err
	}

	if !cfg.Enabled {
		log.Println("AWS integration disabled, skipping sync")
		return nil
	}

	log.Println("Starting AWS SES suppression sync...")

	// 🔥 SES CLIENT DIBUAT DI SINI (BUKAN DI INIT)
	sesClient := aws.NewSESClient(cfg)

	awsSuppressions, err := sesClient.GetSuppressionList(ctx)
	if err != nil {
		log.Printf("AWS failed: %v", err)
		return err
	}

	log.Printf("Retrieved %d suppressions from AWS", len(awsSuppressions))

	if len(awsSuppressions) == 0 {
		log.Println("No suppressions found in AWS")
		return nil
	}

	var suppressions []*models.Suppression
	now := time.Now()

	for _, awsItem := range awsSuppressions {
		suppressions = append(suppressions, &models.Suppression{
			Email:     awsItem.Email,
			Reason:    awsItem.Reason,
			Source:    "AWS",
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	err = s.dbRepo.BulkUpsert(suppressions)
	if err != nil {
		log.Printf("Failed to bulk insert suppressions: %v", err)
		return err
	}

	log.Printf("Sync completed: %d suppressions synced to database", len(suppressions))
	return nil
}

// GetSyncStatus mengembalikan status sync terakhir
func (s *SyncService) GetSyncStatus() (time.Time, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastSync, s.syncInProgress
}
