package handlers

import (
	"net/http"

	"go-backend-valos-id/core/db"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	database *db.Database
}

func NewHealthHandler(database *db.Database) *HealthHandler {
	return &HealthHandler{
		database: database,
	}
}

// Ping handles basic ping endpoint
func (h *HealthHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// HealthCheck checks the health of the application and database
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	status := gin.H{
		"status": "healthy",
	}

	// Check database health
	if err := h.database.Health(); err != nil {
		status["status"] = "unhealthy"
		status["database"] = gin.H{
			"status": "disconnected",
			"error":  err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, status)
		return
	}

	status["database"] = gin.H{
		"status": "connected",
	}

	c.JSON(http.StatusOK, status)
}

// Readiness check for Kubernetes/containers
func (h *HealthHandler) Readiness(c *gin.Context) {
	if err := h.database.Health(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"ready": false,
			"error": "database not ready",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ready": true,
	})
}

// Liveness check for Kubernetes/containers
func (h *HealthHandler) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"alive": true,
	})
}
