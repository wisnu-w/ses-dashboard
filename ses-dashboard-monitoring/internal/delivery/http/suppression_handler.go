package http

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"ses-monitoring/internal/domain/models"
	"ses-monitoring/internal/domain/settings"
	"ses-monitoring/internal/domain/suppression"
	"ses-monitoring/internal/infrastructure/aws"
	"ses-monitoring/internal/infrastructure/database"
	"ses-monitoring/internal/services"

	"github.com/gin-gonic/gin"
)

type SuppressionHandler struct {
	settingsRepo     settings.Repository
	suppressionRepo  suppression.Repository
	suppressionDBRepo *database.SuppressionRepository
	syncService      *services.SyncService
}

func NewSuppressionHandler(settingsRepo settings.Repository, suppressionRepo suppression.Repository, suppressionDBRepo *database.SuppressionRepository, syncService *services.SyncService) *SuppressionHandler {
	return &SuppressionHandler{
		settingsRepo:      settingsRepo,
		suppressionRepo:   suppressionRepo,
		suppressionDBRepo: suppressionDBRepo,
		syncService:       syncService,
	}
}

type SuppressionEntry struct {
	ID               int    `json:"id"`
	Email            string `json:"email"`
	SuppressionType  string `json:"suppression_type"`
	Reason           string `json:"reason"`
	AWSStatus        string `json:"aws_status"`
	IsActive         bool   `json:"is_active"`
	CreatedAt        string `json:"created_at"`
	AddedByName      string `json:"added_by_name,omitempty"`
}

type AddSuppressionRequest struct {
	Email  string `json:"email" binding:"required"`
	Reason string `json:"reason"`
}

type BulkSuppressionRequest struct {
	Emails []string `json:"emails" binding:"required"`
	Reason string   `json:"reason"`
}

type BulkRemoveRequest struct {
	Emails []string `json:"emails" binding:"required"`
}

// GetSuppressions godoc
// @Summary Get suppressions with pagination
// @Description Get list of suppressed emails from suppressions table with pagination
// @Tags suppression
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 50, max: 1000)"
// @Param search query string false "Search term for email, reason, or source"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/suppression [get]
func (h *SuppressionHandler) GetSuppressions(c *gin.Context) {
	// Parse pagination parameters
	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	
	limit := 50 // Default page size
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			if parsed > 1000 {
				parsed = 1000 // Max limit
			}
			limit = parsed
		}
	}
	
	offset := (page - 1) * limit
	searchTerm := c.Query("search")
	
	// Check if AWS integration is enabled
	config, err := h.settingsRepo.GetAWSConfig(c.Request.Context())
	if err != nil || !config.Enabled {
		c.JSON(http.StatusOK, gin.H{
			"suppressions": []interface{}{},
			"total":        0,
			"page":         page,
			"limit":        limit,
			"total_pages":  0,
			"message":      "AWS integration is disabled",
		})
		return
	}
	
	var suppressions []*models.Suppression
	var total int
	
	// Get data with search or without
	if searchTerm != "" {
		suppressions, err = h.suppressionDBRepo.SearchSuppressions(searchTerm, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search suppressions: " + err.Error()})
			return
		}
		total, err = h.suppressionDBRepo.GetSearchCount(searchTerm)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get search count: " + err.Error()})
			return
		}
	} else {
		suppressions, err = h.suppressionDBRepo.GetAll(limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get suppressions: " + err.Error()})
			return
		}
		total, err = h.suppressionDBRepo.GetAllCount()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get total count: " + err.Error()})
			return
		}
	}
	
	totalPages := (total + limit - 1) / limit // Ceiling division
	
	c.JSON(http.StatusOK, gin.H{
		"suppressions": suppressions,
		"total":        total,
		"page":         page,
		"limit":        limit,
		"total_pages":  totalPages,
		"has_next":     page < totalPages,
		"has_prev":     page > 1,
	})
}

// AddSuppression godoc
// @Summary Add email to suppression list with AWS sync
// @Description Add email to local suppression list and sync to AWS SES
// @Tags suppression
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body AddSuppressionRequest true "Email to suppress"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/suppression [post]
func (h *SuppressionHandler) AddSuppression(c *gin.Context) {
	var req AddSuppressionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Get user ID for audit trail
	userID := 0
	if uid, exists := c.Get("user_id"); exists {
		switch v := uid.(type) {
		case int:
			userID = v
		case float64:
			userID = int(v)
		}
	}
	
	// Add to local database first
	entry := &suppression.SuppressionEntry{
		Email:           req.Email,
		SuppressionType: suppression.SuppressionTypeManual,
		Reason:          req.Reason,
		AWSStatus:       suppression.AWSStatusUnknown,
		IsActive:        true,
		AddedBy:         userID,
	}
	
	err := h.suppressionRepo.Add(c.Request.Context(), entry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Try to sync to AWS if enabled
	config, err := h.settingsRepo.GetAWSConfig(c.Request.Context())
	if err == nil && config.Enabled {
		sesClient := aws.NewSESClient(config)
		err = sesClient.AddToSuppression(c.Request.Context(), req.Email, req.Reason)
		if err == nil {
			// Update AWS status if sync successful
			h.suppressionRepo.UpdateAWSStatus(c.Request.Context(), req.Email, suppression.AWSStatusSuppressed)
			h.suppressionRepo.MarkAsSynced(c.Request.Context(), req.Email)
		}
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Email added to suppression list"})
}

// BulkAddSuppression godoc
// @Summary Bulk add emails to suppression list
// @Description Add multiple emails to suppression list at once
// @Tags suppression
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body BulkSuppressionRequest true "Emails to suppress"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /api/suppression/bulk [post]
func (h *SuppressionHandler) BulkAddSuppression(c *gin.Context) {
	var req BulkSuppressionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if len(req.Emails) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one email is required"})
		return
	}
	
	// Get user ID for audit trail
	userID := 0
	if uid, exists := c.Get("user_id"); exists {
		switch v := uid.(type) {
		case int:
			userID = v
		case float64:
			userID = int(v)
		}
	}
	
	successCount := 0
	failedEmails := []string{}
	
	// Get AWS config once
	config, _ := h.settingsRepo.GetAWSConfig(c.Request.Context())
	var sesClient *aws.SESClient
	if config != nil && config.Enabled {
		sesClient = aws.NewSESClient(config)
	}
	
	// Process each email
	for _, email := range req.Emails {
		// Add to local database
		entry := &suppression.SuppressionEntry{
			Email:           email,
			SuppressionType: suppression.SuppressionTypeManual,
			Reason:          req.Reason,
			AWSStatus:       suppression.AWSStatusUnknown,
			IsActive:        true,
			AddedBy:         userID,
		}
		
		err := h.suppressionRepo.Add(c.Request.Context(), entry)
		if err != nil {
			failedEmails = append(failedEmails, email)
			continue
		}
		
		// Try to sync to AWS if enabled
		if sesClient != nil {
			err = sesClient.AddToSuppression(c.Request.Context(), email, req.Reason)
			if err == nil {
				h.suppressionRepo.UpdateAWSStatus(c.Request.Context(), email, suppression.AWSStatusSuppressed)
				h.suppressionRepo.MarkAsSynced(c.Request.Context(), email)
			}
		}
		
		successCount++
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message":       "Bulk suppression completed",
		"success_count": successCount,
		"failed_count":  len(failedEmails),
		"failed_emails": failedEmails,
	})
}

// RemoveSuppression godoc
// @Summary Remove email from AWS SES suppression list
// @Description Remove email from AWS SES suppression list
// @Tags suppression
// @Produce json
// @Security BearerAuth
// @Param email path string true "Email address"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/suppression/{email} [delete]
func (h *SuppressionHandler) RemoveSuppression(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}
	
	config, err := h.settingsRepo.GetAWSConfig(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	if !config.Enabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "AWS integration is disabled"})
		return
	}
	
	sesClient := aws.NewSESClient(config)
	err = sesClient.RemoveFromSuppression(c.Request.Context(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// If AWS removal successful, also remove from local DB
	if err := h.suppressionDBRepo.Delete(email); err != nil {
		log.Printf("Warning: Failed to remove %s from local DB: %v", email, err)
		// Don't fail the operation, just log the warning
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Email removed from AWS SES suppression list and local database"})
}

// BulkRemoveSuppression godoc
// @Summary Bulk remove emails from suppression list
// @Description Remove multiple emails from AWS SES suppression list at once
// @Tags suppression
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body BulkRemoveRequest true "Emails to remove from suppression"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /api/suppression/bulk [delete]
func (h *SuppressionHandler) BulkRemoveSuppression(c *gin.Context) {
	var req BulkRemoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if len(req.Emails) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one email is required"})
		return
	}
	
	config, err := h.settingsRepo.GetAWSConfig(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	if !config.Enabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "AWS integration is disabled"})
		return
	}
	
	sesClient := aws.NewSESClient(config)
	successCount := 0
	failedEmails := []string{}
	
	// Process each email
	for _, email := range req.Emails {
		err := sesClient.RemoveFromSuppression(c.Request.Context(), email)
		if err != nil {
			log.Printf("Failed to remove %s from AWS: %v", email, err)
			failedEmails = append(failedEmails, email)
			continue
		}
		
		// If AWS removal successful, also remove from local DB
		if err := h.suppressionDBRepo.Delete(email); err != nil {
			log.Printf("Warning: Failed to remove %s from local DB: %v", email, err)
			// Don't fail the operation, just log the warning
		}
		
		log.Printf("Successfully removed %s from AWS and local DB", email)
		successCount++
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message":       "Bulk removal completed",
		"success_count": successCount,
		"failed_count":  len(failedEmails),
		"failed_emails": failedEmails,
		"details":       "Check server logs for detailed error messages",
	})
}

// SyncFromAWS godoc
// @Summary Trigger manual sync from AWS SES
// @Description Trigger manual sync of suppression list from AWS SES to local database
// @Tags suppression
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/suppression/sync [post]
func (h *SuppressionHandler) SyncFromAWS(c *gin.Context) {
	config, err := h.settingsRepo.GetAWSConfig(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get AWS config: " + err.Error()})
		return
	}
	
	if !config.Enabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "AWS integration is disabled. Please configure AWS settings first."})
		return
	}
	
	// Trigger manual sync dengan context baru (tidak terikat request)
	go h.syncService.SyncNow(context.Background())
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Sync triggered. Data will be updated in background.",
		"status":  "in_progress",
	})
}

// GetSyncStatus godoc
// @Summary Get sync status
// @Description Get last sync time and current sync status
// @Tags suppression
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/suppression/sync/status [get]
func (h *SuppressionHandler) GetSyncStatus(c *gin.Context) {
	lastSync, inProgress := h.syncService.GetSyncStatus()
	
	// Check if AWS integration is enabled
	config, err := h.settingsRepo.GetAWSConfig(c.Request.Context())
	count := 0
	if err == nil && config.Enabled {
		// Use count method instead of loading all data
		if dbCount, err := h.suppressionDBRepo.CountBySource("AWS"); err == nil {
			count = dbCount
		}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"last_sync":    lastSync,
		"in_progress":  inProgress,
		"next_sync_in": "5 minutes",
		"db_count":     count,
		"aws_enabled":  config != nil && config.Enabled,
	})
}