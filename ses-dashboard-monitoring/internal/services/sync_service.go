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

// StartBackgroundSync memulai sync otomatis dengan interval dinamis
func (s *SyncService) StartBackgroundSync(ctx context.Context) {
	// Sync pertama setelah startup
	go func() {
		time.Sleep(10 * time.Second)
		_ = s.SyncNow(context.Background())
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Get current sync interval from settings
			cfg, err := s.settingsRepo.GetAWSConfig(context.Background())
			if err != nil || !cfg.Enabled {
				time.Sleep(5 * time.Minute) // fallback interval
				continue
			}
			
			interval := time.Duration(cfg.SyncInterval) * time.Minute
			if interval < time.Minute {
				interval = 5 * time.Minute // minimum 1 minute
			}
			
			time.Sleep(interval)
			go s.SyncNow(context.Background())
		}
	}
}

// SyncNow melakukan sync manual dengan full synchronization
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

	// Create context with timeout for the entire sync operation
	syncCtx, cancel := context.WithTimeout(ctx, 25*time.Minute)
	defer cancel()

	// 🔥 AMBIL CONFIG TERBARU DARI DB
	cfg, err := s.settingsRepo.GetAWSConfig(syncCtx)
	if err != nil {
		log.Printf("Failed to load AWS config: %v", err)
		return err
	}

	if !cfg.Enabled {
		log.Println("AWS integration disabled, skipping sync")
		return nil
	}

	if cfg.AccessKey == "" || cfg.SecretKey == "" {
		log.Println("AWS credentials not configured, skipping sync")
		return nil
	}

	log.Printf("Starting AWS SES suppression sync (region: %s)...", cfg.Region)

	// 🔥 SES CLIENT DIBUAT DI SINI (BUKAN DI INIT)
	sesClient := aws.NewSESClient(cfg)

	// Test connection first
	if err := sesClient.TestConnection(syncCtx); err != nil {
		log.Printf("AWS connection test failed: %v", err)
		return err
	}

	// Get all AWS suppressions
	awsSuppressions, err := sesClient.GetSuppressionList(syncCtx)
	if err != nil {
		log.Printf("AWS sync failed: %v", err)
		return err
	}

	log.Printf("Retrieved %d suppressions from AWS", len(awsSuppressions))

	// Get all current DB suppressions from AWS source
	dbSuppressions, err := s.dbRepo.GetBySource("AWS")
	if err != nil {
		log.Printf("Failed to get DB suppressions: %v", err)
		return err
	}

	log.Printf("Found %d AWS suppressions in database", len(dbSuppressions))

	// Create maps for comparison
	awsMap := make(map[string]*aws.SuppressionStatus)
	dbMap := make(map[string]*models.Suppression)

	for _, item := range awsSuppressions {
		awsMap[item.Email] = item
	}

	for _, item := range dbSuppressions {
		dbMap[item.Email] = item
	}

	// Find emails to add (in AWS but not in DB)
	var toAdd []*models.Suppression
	for email, awsItem := range awsMap {
		if _, exists := dbMap[email]; !exists {
			toAdd = append(toAdd, &models.Suppression{
				Email:     awsItem.Email,
				Reason:    awsItem.Reason,
				Source:    "AWS",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		}
	}

	// Find emails to remove (in DB but not in AWS)
	var toRemove []string
	for email := range dbMap {
		if _, exists := awsMap[email]; !exists {
			toRemove = append(toRemove, email)
		}
	}

	// Find emails to update (different reason)
	var toUpdate []*models.Suppression
	for email, awsItem := range awsMap {
		if dbItem, exists := dbMap[email]; exists {
			if dbItem.Reason != awsItem.Reason {
				toUpdate = append(toUpdate, &models.Suppression{
					Email:     awsItem.Email,
					Reason:    awsItem.Reason,
					Source:    "AWS",
					CreatedAt: dbItem.CreatedAt,
					UpdatedAt: time.Now(),
				})
			}
		}
	}

	log.Printf("Sync plan: %d to add, %d to remove, %d to update", len(toAdd), len(toRemove), len(toUpdate))

	// Execute sync operations
	if len(toAdd) > 0 {
		if err := s.dbRepo.BulkUpsert(toAdd); err != nil {
			log.Printf("Failed to add suppressions: %v", err)
			return err
		}
		log.Printf("Added %d suppressions", len(toAdd))
	}

	if len(toUpdate) > 0 {
		if err := s.dbRepo.BulkUpsert(toUpdate); err != nil {
			log.Printf("Failed to update suppressions: %v", err)
			return err
		}
		log.Printf("Updated %d suppressions", len(toUpdate))
	}

	if len(toRemove) > 0 {
		if err := s.dbRepo.BulkDelete(toRemove); err != nil {
			log.Printf("Failed to remove suppressions: %v", err)
			return err
		}
		log.Printf("Removed %d suppressions", len(toRemove))
	}

	// Final count verification
	finalCount, err := s.dbRepo.CountBySource("AWS")
	if err != nil {
		log.Printf("Failed to get final count: %v", err)
	} else {
		log.Printf("Sync completed successfully: AWS=%d, DB=%d", len(awsSuppressions), finalCount)
		if finalCount != len(awsSuppressions) {
			log.Printf("WARNING: Count mismatch after sync! AWS=%d, DB=%d", len(awsSuppressions), finalCount)
		}
	}

	return nil
}

// GetSyncStatus mengembalikan status sync terakhir
func (s *SyncService) GetSyncStatus() (time.Time, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastSync, s.syncInProgress
}
