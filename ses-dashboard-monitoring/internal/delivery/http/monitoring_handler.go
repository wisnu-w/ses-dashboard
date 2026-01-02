package http

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"ses-monitoring/internal/domain/sesevent"
	"ses-monitoring/internal/domain/settings"
	"ses-monitoring/internal/usecase"

	"github.com/gin-gonic/gin"
)

type MonitoringHandler struct {
	uc           *usecase.SESUsecase
	settingsRepo settings.Repository

	// Timezone cache
	timezoneMu    sync.RWMutex
	timezoneCache string
}

func NewMonitoringHandler(uc *usecase.SESUsecase, settingsRepo settings.Repository) *MonitoringHandler {
	h := &MonitoringHandler{
		uc:            uc,
		settingsRepo:  settingsRepo,
		timezoneCache: "Asia/Jakarta", // default
	}
	h.loadTimezone() // Load initial timezone
	return h
}

// GetEvents godoc
// @Summary Get SES events with pagination
// @Description Retrieve list of SES events with pagination support
// @Tags monitoring
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)" minimum(1)
// @Param limit query int false "Number of events per page (default: 50, max: 1000)" minimum(1) maximum(1000)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/events [get]
func (h *MonitoringHandler) GetEvents(c *gin.Context) {
	// Parse query parameters
	page := 1
	limit := 50
	search := c.Query("search")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 1000 {
			limit = parsed
		}
	}

	offset := (page - 1) * limit

	var events []*sesevent.Event
	var total int
	var err error

	// Use optimized queries based on filter presence
	if search != "" || startDate != "" || endDate != "" {
		events, err = h.uc.GetEventsWithFilter(c.Request.Context(), limit, offset, search, startDate, endDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Get filtered count
		total, err = h.uc.GetFilteredEventCount(c.Request.Context(), search, startDate, endDate)
	} else {
		events, err = h.uc.GetEventsPaginated(c.Request.Context(), limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Get total count
		total, err = h.uc.GetEventCount(c.Request.Context())
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (total + limit - 1) / limit // Ceiling division

	c.JSON(http.StatusOK, gin.H{
		"events": events,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": totalPages,
			"hasNext":    page < totalPages,
			"hasPrev":    page > 1,
		},
	})
}

type MetricsResponse struct {
	TotalEvents    int     `json:"total_events"`
	SendCount      int     `json:"send_count"`
	DeliveryCount  int     `json:"delivery_count"`
	BounceCount    int     `json:"bounce_count"`
	ComplaintCount int     `json:"complaint_count"`
	OpenCount      int     `json:"open_count"`
	ClickCount     int     `json:"click_count"`
	BounceRate     float64 `json:"bounce_rate"`
	DeliveryRate   float64 `json:"delivery_rate"`
}

// GetMetrics godoc
// @Summary Get overall metrics
// @Description Retrieve overall SES metrics with counts and rates
// @Tags monitoring
// @Produce json
// @Security BearerAuth
// @Success 200 {object} MetricsResponse
// @Failure 500 {object} map[string]string
// @Router /api/metrics [get]
func (h *MonitoringHandler) GetMetrics(c *gin.Context) {
	total, err := h.uc.GetEventCount(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bounceRate, err := h.uc.GetBounceRate(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	deliveryRate, err := h.uc.GetDeliveryRate(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	counts, err := h.uc.GetEventTypeCounts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	metrics := MetricsResponse{
		TotalEvents:    total,
		SendCount:      counts["Send"],
		DeliveryCount:  counts["Delivery"],
		BounceCount:    counts["Bounce"],
		ComplaintCount: counts["Complaint"],
		OpenCount:      counts["Open"],
		ClickCount:     counts["Click"],
		BounceRate:     bounceRate,
		DeliveryRate:   deliveryRate,
	}

	c.JSON(http.StatusOK, metrics)
}

// GetDailyMetrics godoc
// @Summary Get daily metrics
// @Description Retrieve daily SES metrics
// @Tags monitoring
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string][]sesevent.DailyMetrics
// @Failure 500 {object} map[string]string
// @Router /api/metrics/daily [get]
func (h *MonitoringHandler) GetDailyMetrics(c *gin.Context) {
	metrics, err := h.uc.GetDailyMetrics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert timestamps to configured timezone
	if err := h.convertMetricsTimezone(metrics); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"daily_metrics": metrics})
}

// GetMonthlyMetrics godoc
// @Summary Get monthly metrics
// @Description Retrieve monthly SES metrics
// @Tags monitoring
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string][]sesevent.MonthlyMetrics
// @Failure 500 {object} map[string]string
// @Router /api/metrics/monthly [get]
func (h *MonitoringHandler) GetMonthlyMetrics(c *gin.Context) {
	metrics, err := h.uc.GetMonthlyMetrics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert timestamps to configured timezone
	if err := h.convertMetricsTimezone(metrics); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"monthly_metrics": metrics})
}

// GetHourlyMetrics godoc
// @Summary Get hourly metrics
// @Description Retrieve hourly SES metrics
// @Tags monitoring
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string][]sesevent.HourlyMetrics
// @Failure 500 {object} map[string]string
// @Router /api/metrics/hourly [get]
func (h *MonitoringHandler) GetHourlyMetrics(c *gin.Context) {
	metrics, err := h.uc.GetHourlyMetrics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert timestamps to configured timezone
	if err := h.convertMetricsTimezone(metrics); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"hourly_metrics": metrics})
}

func (h *MonitoringHandler) convertMetricsTimezone(metrics interface{}) error {
	timezone := h.getTimezoneFromCache()

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return err
	}

	switch m := metrics.(type) {
	case []sesevent.DailyMetrics:
		for i := range m {
			if t, err := time.Parse("2006-01-02", m[i].Date); err == nil {
				m[i].Date = t.In(loc).Format("2006-01-02")
			}
		}
	case []sesevent.MonthlyMetrics:
		for i := range m {
			if t, err := time.Parse("2006-01", m[i].Month); err == nil {
				m[i].Month = t.In(loc).Format("2006-01")
			}
		}
	case []sesevent.HourlyMetrics:
		for i := range m {
			if t, err := time.Parse("2006-01-02 15:04", m[i].Hour); err == nil {
				m[i].Hour = t.In(loc).Format("2006-01-02 15:04")
			}
		}
	}

	return nil
}

func (h *MonitoringHandler) getTimezoneFromCache() string {
	h.timezoneMu.RLock()
	timezone := h.timezoneCache
	h.timezoneMu.RUnlock()
	return timezone
}

func (h *MonitoringHandler) loadTimezone() {
	if config, err := h.settingsRepo.GetTimezoneConfig(context.Background()); err == nil {
		h.timezoneMu.Lock()
		h.timezoneCache = config.Timezone
		h.timezoneMu.Unlock()
	}
}

func (h *MonitoringHandler) RefreshTimezone() {
	h.loadTimezone()
}
