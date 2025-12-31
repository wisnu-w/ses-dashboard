package services

import (
	"context"
	"log"
	"sync"
	"time"

	"ses-monitoring/internal/domain/models"
	"ses-monitoring/internal/infrastructure/aws"
	"ses-monitoring/internal/infrastructure/database"
)

type SyncService struct {
	awsClient *aws.SESClient
	dbRepo    *database.SuppressionRepository
	mu        sync.RWMutex
	lastSync  time.Time
	syncInProgress bool
}

func NewSyncService(awsClient *aws.SESClient, dbRepo *database.SuppressionRepository) *SyncService {
	return &SyncService{
		awsClient: awsClient,
		dbRepo:    dbRepo,
	}
}

// StartBackgroundSync memulai sync otomatis setiap 5 menit (untuk testing)
func (s *SyncService) StartBackgroundSync(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute) // Lebih sering untuk testing
	defer ticker.Stop()

	// Sync pertama kali saat startup (dengan delay 10 detik)
	go func() {
		time.Sleep(10 * time.Second)
		s.SyncNow(ctx)
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			go s.SyncNow(ctx)
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

	log.Println("Starting AWS SES suppression sync...")

	// Get all data dari AWS
	awsSuppressions, err := s.awsClient.GetSuppressionList(ctx)
	if err != nil {
		log.Printf("AWS failed: %v", err)
		return err
	}

	log.Printf("Retrieved %d suppressions from AWS", len(awsSuppressions))

	if len(awsSuppressions) == 0 {
		log.Println("No suppressions found in AWS")
		return nil
	}

	// Convert semua ke domain models
	var suppressions []*models.Suppression
	for _, aws := range awsSuppressions {
		suppressions = append(suppressions, &models.Suppression{
			Email:     aws.Email,
			Reason:    aws.Reason,
			Source:    "AWS",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	// Bulk insert semua sekaligus
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