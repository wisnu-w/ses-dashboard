package http

import (
	"context"
	"errors"
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

	metricsCacheMu sync.RWMutex
	metricsCache   map[string]metricsCacheItem
}

type metricsCacheItem struct {
	data      interface{}
	expiresAt time.Time
}

const metricsCacheTTL = 30 * time.Second

func NewMonitoringHandler(uc *usecase.SESUsecase, settingsRepo settings.Repository) *MonitoringHandler {
	h := &MonitoringHandler{
		uc:            uc,
		settingsRepo:  settingsRepo,
		timezoneCache: "Asia/Jakarta", // default
		metricsCache:  make(map[string]metricsCacheItem),
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

	// Convert timestamps to configured timezone
	if err := h.convertEventsTimezone(events); err != nil {
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
	cacheKey := "summary:" + h.getTimezoneFromCache()
	if cached, ok := h.getMetricsCache(cacheKey); ok {
		if metrics, ok := cached.(MetricsResponse); ok {
			c.JSON(http.StatusOK, metrics)
			return
		}
	}

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

	h.setMetricsCache(cacheKey, metrics)
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
	now := time.Now().UTC()
	start, end, err := h.parseDateRange(
		c,
		now.AddDate(0, 0, -30),
		now,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cacheKey := h.buildMetricsCacheKey("daily", start, end)
	if cached, ok := h.getMetricsCache(cacheKey); ok {
		if metrics, ok := cached.([]*sesevent.DailyMetrics); ok {
			c.JSON(http.StatusOK, gin.H{"daily_metrics": metrics})
			return
		}
	}

	metrics, err := h.uc.GetDailyMetrics(c.Request.Context(), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert timestamps to configured timezone
	if err := h.convertMetricsTimezone(metrics); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.setMetricsCache(cacheKey, metrics)
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
	now := time.Now().UTC()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	start, end, err := h.parseDateRange(
		c,
		startOfMonth.AddDate(0, -11, 0),
		now,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cacheKey := h.buildMetricsCacheKey("monthly", start, end)
	if cached, ok := h.getMetricsCache(cacheKey); ok {
		if metrics, ok := cached.([]*sesevent.MonthlyMetrics); ok {
			c.JSON(http.StatusOK, gin.H{"monthly_metrics": metrics})
			return
		}
	}

	metrics, err := h.uc.GetMonthlyMetrics(c.Request.Context(), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert timestamps to configured timezone
	if err := h.convertMetricsTimezone(metrics); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.setMetricsCache(cacheKey, metrics)
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
	now := time.Now().UTC()
	start, end, err := h.parseDateRange(
		c,
		now.Add(-48*time.Hour),
		now,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cacheKey := h.buildMetricsCacheKey("hourly", start, end)
	if cached, ok := h.getMetricsCache(cacheKey); ok {
		if metrics, ok := cached.([]*sesevent.HourlyMetrics); ok {
			c.JSON(http.StatusOK, gin.H{"hourly_metrics": metrics})
			return
		}
	}

	metrics, err := h.uc.GetHourlyMetrics(c.Request.Context(), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert timestamps to configured timezone
	if err := h.convertMetricsTimezone(metrics); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.setMetricsCache(cacheKey, metrics)
	c.JSON(http.StatusOK, gin.H{"hourly_metrics": metrics})
}

func (h *MonitoringHandler) buildMetricsCacheKey(prefix string, start, end *time.Time) string {
	timezone := h.getTimezoneFromCache()
	if start == nil && end == nil {
		return prefix + ":none:" + timezone
	}
	startVal := ""
	endVal := ""
	if start != nil {
		startVal = start.Format(time.RFC3339)
	}
	if end != nil {
		endVal = end.Format(time.RFC3339)
	}
	return prefix + ":" + startVal + ":" + endVal + ":" + timezone
}

func (h *MonitoringHandler) getMetricsCache(key string) (interface{}, bool) {
	h.metricsCacheMu.RLock()
	item, ok := h.metricsCache[key]
	h.metricsCacheMu.RUnlock()
	if !ok {
		return nil, false
	}
	if time.Now().After(item.expiresAt) {
		h.metricsCacheMu.Lock()
		delete(h.metricsCache, key)
		h.metricsCacheMu.Unlock()
		return nil, false
	}
	return item.data, true
}

func (h *MonitoringHandler) setMetricsCache(key string, data interface{}) {
	h.metricsCacheMu.Lock()
	h.metricsCache[key] = metricsCacheItem{
		data:      data,
		expiresAt: time.Now().Add(metricsCacheTTL),
	}
	h.metricsCacheMu.Unlock()
}

func (h *MonitoringHandler) parseDateRange(c *gin.Context, defaultStart, defaultEnd time.Time) (*time.Time, *time.Time, error) {
	start := defaultStart
	end := defaultEnd

	if startQuery := c.Query("start_date"); startQuery != "" {
		parsed, err := time.Parse("2006-01-02", startQuery)
		if err != nil {
			return nil, nil, err
		}
		start = parsed.UTC()
	}

	if endQuery := c.Query("end_date"); endQuery != "" {
		parsed, err := time.Parse("2006-01-02", endQuery)
		if err != nil {
			return nil, nil, err
		}
		end = parsed.UTC().Add(24 * time.Hour)
	}

	if start.After(end) {
		return nil, nil, errors.New("start_date must be before end_date")
	}

	return &start, &end, nil
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
				utcTime := t.UTC()
				m[i].Date = utcTime.In(loc).Format("2006-01-02")
			}
		}
	case []*sesevent.DailyMetrics:
		for i := range m {
			if t, err := time.Parse("2006-01-02", m[i].Date); err == nil {
				utcTime := t.UTC()
				m[i].Date = utcTime.In(loc).Format("2006-01-02")
			}
		}
	case []sesevent.MonthlyMetrics:
		for i := range m {
			if t, err := time.Parse("2006-01", m[i].Month); err == nil {
				utcTime := t.UTC()
				m[i].Month = utcTime.In(loc).Format("2006-01")
			}
		}
	case []*sesevent.MonthlyMetrics:
		for i := range m {
			if t, err := time.Parse("2006-01", m[i].Month); err == nil {
				utcTime := t.UTC()
				m[i].Month = utcTime.In(loc).Format("2006-01")
			}
		}
	case []sesevent.HourlyMetrics:
		for i := range m {
			if t, err := time.Parse("2006-01-02 15:04", m[i].Hour); err == nil {
				// Assume database time is UTC
				utcTime := t.UTC()
				m[i].Hour = utcTime.In(loc).Format("2006-01-02 15:04")
			}
		}
	case []*sesevent.HourlyMetrics:
		for i := range m {
			if t, err := time.Parse("2006-01-02 15:04", m[i].Hour); err == nil {
				// Assume database time is UTC
				utcTime := t.UTC()
				m[i].Hour = utcTime.In(loc).Format("2006-01-02 15:04")
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

func (h *MonitoringHandler) convertEventsTimezone(events []*sesevent.Event) error {
	timezone := h.getTimezoneFromCache()

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return err
	}

	for i := range events {
		// Convert EventTimestamp
		events[i].EventTimestamp = events[i].EventTimestamp.In(loc)
		// Convert CreatedAt
		events[i].CreatedAt = events[i].CreatedAt.In(loc)
	}

	return nil
}
