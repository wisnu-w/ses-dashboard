package http

import (
	"net/http"
	"strconv"

	"ses-monitoring/internal/domain/settings"
	"ses-monitoring/internal/infrastructure/aws"

	"github.com/gin-gonic/gin"
)

type SettingsHandler struct {
	settingsRepo settings.Repository
}

func NewSettingsHandler(settingsRepo settings.Repository) *SettingsHandler {
	return &SettingsHandler{settingsRepo: settingsRepo}
}

// GetAWSSettings godoc
// @Summary Get AWS settings
// @Description Get current AWS configuration settings
// @Tags settings
// @Produce json
// @Security BearerAuth
// @Success 200 {object} settings.AWSConfig
// @Failure 500 {object} map[string]string
// @Router /api/settings/aws [get]
func (h *SettingsHandler) GetAWSSettings(c *gin.Context) {
	config, err := h.settingsRepo.GetAWSConfig(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Don't return secret key in response
	config.SecretKey = ""
	if config.AccessKey != "" && len(config.AccessKey) > 4 {
		config.AccessKey = config.AccessKey[:4] + "****"
	} else if config.AccessKey != "" {
		config.AccessKey = "****"
	}
	
	c.JSON(http.StatusOK, config)
}

// UpdateAWSSettings godoc
// @Summary Update AWS settings
// @Description Update AWS configuration settings
// @Tags settings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param settings body settings.AWSConfig true "AWS settings"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/settings/aws [put]
func (h *SettingsHandler) UpdateAWSSettings(c *gin.Context) {
	var config settings.AWSConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Get user ID from JWT token - required for audit trail
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User authentication required"})
		return
	}
	
	// Handle both int and float64 from JWT claims
	var userIDInt int
	switch v := userID.(type) {
	case int:
		userIDInt = v
	case float64:
		userIDInt = int(v)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
		return
	}
	
	if userIDInt == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}
	
	// Save settings with proper user tracking
	ctx := c.Request.Context()
	
	err := h.settingsRepo.Set(ctx, "aws_enabled", strconv.FormatBool(config.Enabled), userIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	err = h.settingsRepo.Set(ctx, "aws_region", config.Region, userIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	if config.AccessKey != "" {
		err = h.settingsRepo.Set(ctx, "aws_access_key", config.AccessKey, userIDInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	
	if config.SecretKey != "" {
		err = h.settingsRepo.Set(ctx, "aws_secret_key", config.SecretKey, userIDInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Settings updated successfully"})
}

// TestAWSConnection godoc
// @Summary Test AWS connection
// @Description Test AWS SES connection with provided credentials
// @Tags settings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param settings body settings.AWSConfig true "AWS settings to test"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/settings/aws/test [post]
func (h *SettingsHandler) TestAWSConnection(c *gin.Context) {
	var config settings.AWSConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	sesClient := aws.NewSESClient(&config)
	err := sesClient.TestConnection(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "AWS connection successful"})
}

type RetentionSettings struct {
	RetentionDays int  `json:"retention_days"` // 0 = never delete
	Enabled       bool `json:"enabled"`
}

// GetRetentionSettings godoc
// @Summary Get retention settings
// @Description Get event log retention settings
// @Tags settings
// @Produce json
// @Security BearerAuth
// @Success 200 {object} RetentionSettings
// @Router /api/settings/retention [get]
func (h *SettingsHandler) GetRetentionSettings(c *gin.Context) {
	retentionDaysSetting, err1 := h.settingsRepo.Get(c.Request.Context(), "retention_days")
	retentionEnabledSetting, err2 := h.settingsRepo.Get(c.Request.Context(), "retention_enabled")
	
	// Debug logging
	if err1 != nil {
		c.Header("X-Debug-Days-Error", err1.Error())
	} else {
		c.Header("X-Debug-Days-Value", retentionDaysSetting.Value)
	}
	
	if err2 != nil {
		c.Header("X-Debug-Enabled-Error", err2.Error())
	} else {
		c.Header("X-Debug-Enabled-Value", retentionEnabledSetting.Value)
	}
	
	var days int
	if err1 == nil {
		if parsed, err := strconv.Atoi(retentionDaysSetting.Value); err == nil {
			days = parsed
		}
	}
	
	var enabled bool
	if err2 == nil {
		enabled = retentionEnabledSetting.Value == "true"
	}
	
	c.JSON(http.StatusOK, RetentionSettings{
		RetentionDays: days,
		Enabled:       enabled,
	})
}

// UpdateRetentionSettings godoc
// @Summary Update retention settings
// @Description Update event log retention settings
// @Tags settings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param settings body RetentionSettings true "Retention settings"
// @Success 200 {object} map[string]string
// @Router /api/settings/retention [put]
func (h *SettingsHandler) UpdateRetentionSettings(c *gin.Context) {
	var settings RetentionSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User authentication required"})
		return
	}
	
	var userIDInt int
	switch v := userID.(type) {
	case int:
		userIDInt = v
	case float64:
		userIDInt = int(v)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
		return
	}
	
	ctx := c.Request.Context()
	
	err := h.settingsRepo.Set(ctx, "retention_days", strconv.Itoa(settings.RetentionDays), userIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	err = h.settingsRepo.Set(ctx, "retention_enabled", strconv.FormatBool(settings.Enabled), userIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Retention settings updated successfully"})
}
// @Summary Check email suppression status
// @Description Check if email is suppressed in AWS SES
// @Tags suppression
// @Produce json
// @Security BearerAuth
// @Param email path string true "Email address"
// @Success 200 {object} aws.SuppressionStatus
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/suppression/{email}/status [get]
func (h *SettingsHandler) CheckEmailSuppression(c *gin.Context) {
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
	
	sesClient := aws.NewSESClient(config)
	status, err := sesClient.CheckSuppressionStatus(c.Request.Context(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, status)
}

// RemoveEmailSuppression godoc
// @Summary Remove email from suppression
// @Description Remove email from AWS SES suppression list
// @Tags suppression
// @Produce json
// @Security BearerAuth
// @Param email path string true "Email address"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/suppression/{email} [delete]
func (h *SettingsHandler) RemoveEmailSuppression(c *gin.Context) {
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
	
	sesClient := aws.NewSESClient(config)
	err = sesClient.RemoveFromSuppression(c.Request.Context(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Email removed from suppression list"})
}