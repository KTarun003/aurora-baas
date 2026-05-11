package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	startTime time.Time
}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
	}
}

// Check returns the health status of the service
// @Summary Check service health
// @Description Get the health status and uptime of the Aurora service
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	uptime := int64(time.Since(h.startTime).Seconds())

	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"uptime":  uptime,
		"service": "aurora-core",
	})
}
