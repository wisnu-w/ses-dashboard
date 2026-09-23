package services

import (
	"context"
	"log"
	"sync"
	"time"

	"ses-monitoring/internal/domain/sesevent"
)

type MVRefreshService struct {
	sesRepo    sesevent.Repository
	mu         sync.Mutex
	refreshing bool
}

func NewMVRefreshService(sesRepo sesevent.Repository) *MVRefreshService {
	return &MVRefreshService{
		sesRepo: sesRepo,
	}
}

// StartRefreshScheduler memulai refresh Materialized View secara periodik
func (s *MVRefreshService) StartRefreshScheduler(ctx context.Context) {
	// Refresh setiap 5 menit
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	// Refresh pertama kali setelah 30 detik startup (memberikan waktu bagi DB dll)
	go func() {
		time.Sleep(30 * time.Second)
		s.RunRefresh(ctx)
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			go s.RunRefresh(ctx)
		}
	}
}

// RunRefresh menjalankan perintah REFRESH MATERIALIZED VIEW CONCURRENTLY
func (s *MVRefreshService) RunRefresh(ctx context.Context) {
	s.mu.Lock()
	if s.refreshing {
		s.mu.Unlock()
		log.Println("Refresh MV already in progress, skipping...")
		return
	}
	s.refreshing = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.refreshing = false
		s.mu.Unlock()
	}()

	log.Println("Starting Materialized View refresh (mv_ses_daily_summary)...")

	// Set timeout agar tidak nge-hang jika DB lock terlalu lama
	refreshCtx, cancel := context.WithTimeout(ctx, 4*time.Minute)
	defer cancel()

	err := s.sesRepo.RefreshDailySummary(refreshCtx)
	if err != nil {
		log.Printf("Failed to refresh mv_ses_daily_summary: %v", err)
		return
	}

	log.Println("Successfully refreshed mv_ses_daily_summary")
}
