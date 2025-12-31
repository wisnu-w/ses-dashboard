package services

import (
	"context"
	"log"
	"strconv"
	"time"

	"ses-monitoring/internal/domain/sesevent"
	"ses-monitoring/internal/domain/settings"
)

type CleanupService struct {
	settingsRepo settings.Repository
	sesRepo      sesevent.Repository
}

func NewCleanupService(settingsRepo settings.Repository, sesRepo sesevent.Repository) *CleanupService {
	return &CleanupService{
		settingsRepo: settingsRepo,
		sesRepo:      sesRepo,
	}
}

// StartCleanupScheduler memulai cleanup otomatis setiap hari
func (s *CleanupService) StartCleanupScheduler(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour) // Cleanup setiap hari
	defer ticker.Stop()

	// Cleanup pertama kali setelah 1 menit startup
	go func() {
		time.Sleep(1 * time.Minute)
		s.RunCleanup(ctx)
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			go s.RunCleanup(ctx)
		}
	}
}

// RunCleanup menjalankan cleanup berdasarkan retention settings
func (s *CleanupService) RunCleanup(ctx context.Context) {
	log.Println("Starting event log cleanup...")

	// Get retention settings
	retentionEnabledSetting, err := s.settingsRepo.Get(ctx, "retention_enabled")
	if err != nil || retentionEnabledSetting == nil || retentionEnabledSetting.Value != "true" {
		log.Println("Retention cleanup is disabled")
		return
	}

	retentionDaysSetting, err := s.settingsRepo.Get(ctx, "retention_days")
	if err != nil || retentionDaysSetting == nil {
		log.Println("Retention days not configured")
		return
	}

	retentionDays := 30 // default
	if retentionDaysSetting.Value != "" {
		if parsed, err := strconv.Atoi(retentionDaysSetting.Value); err == nil {
			retentionDays = parsed
		}
	}

	// Skip cleanup if retention_days = 0 (never delete)
	if retentionDays == 0 {
		log.Println("Retention set to never delete - skipping cleanup")
		return
	}

	// Calculate cutoff date
	cutoffDate := time.Now().UTC().AddDate(0, 0, -retentionDays)
	
	log.Printf("Deleting events older than %d days (before %s)", retentionDays, cutoffDate.Format("2006-01-02"))

	// Delete old events
	deletedCount, err := s.sesRepo.DeleteOldEvents(ctx, cutoffDate)
	if err != nil {
		log.Printf("Failed to delete old events: %v", err)
		return
	}

	log.Printf("Cleanup completed: %d old events deleted", deletedCount)
}