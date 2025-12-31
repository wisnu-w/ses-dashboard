package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Service   string    `json:"service"`
}

// Health godoc
// @Summary Health check endpoint
// @Description Check if the service is running
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Service:   "SES Monitoring API",
	}
	
	c.JSON(http.StatusOK, response)
}

// Ready godoc
// @Summary Readiness check endpoint
// @Description Check if the service is ready to serve requests
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /ready [get]
func (h *HealthHandler) Ready(c *gin.Context) {
	response := HealthResponse{
		Status:    "ready",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Service:   "SES Monitoring API",
	}
	
	c.JSON(http.StatusOK, response)
}