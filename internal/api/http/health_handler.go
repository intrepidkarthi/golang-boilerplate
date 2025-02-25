// Package http provides HTTP handlers for the application.
//
// The health handler implements health check endpoints that provide information
// about the service's health and its dependencies (database, cache, etc.).
package http

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"go-boilerplate/internal/cache"
)

// HealthHandler handles health check related endpoints
type HealthHandler struct {
	db    *pgxpool.Pool
	cache *cache.RedisCache
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *pgxpool.Pool, cache *cache.RedisCache) *HealthHandler {
	return &HealthHandler{
		db:    db,
		cache: cache,
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string           `json:"version"`
	Services  map[string]Status `json:"services"`
}

// Status represents the status of a service
type Status struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// Health godoc
// @Summary Get service health status
// @Description Get the health status of the service and its dependencies
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *HealthHandler) Health(c echo.Context) error {
	ctx := context.Background()
	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Version:   "1.0.0", // This should come from build info
		Services:  make(map[string]Status),
	}

	// Check database
	if err := h.db.Ping(ctx); err != nil {
		response.Status = "degraded"
		response.Services["database"] = Status{
			Status:  "down",
			Message: "Database connection failed",
		}
	} else {
		response.Services["database"] = Status{
			Status: "up",
		}
	}

	// Check Redis
	if err := h.cache.Client().Ping(ctx).Err(); err != nil {
		response.Status = "degraded"
		response.Services["cache"] = Status{
			Status:  "down",
			Message: "Redis connection failed",
		}
	} else {
		response.Services["cache"] = Status{
			Status: "up",
		}
	}

	return c.JSON(http.StatusOK, response)
}

// LivenessProbe godoc
// @Summary Kubernetes liveness probe
// @Description Endpoint for Kubernetes liveness probe
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health/live [get]
func (h *HealthHandler) LivenessProbe(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "alive",
	})
}

// ReadinessProbe godoc
// @Summary Kubernetes readiness probe
// @Description Endpoint for Kubernetes readiness probe
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health/ready [get]
func (h *HealthHandler) ReadinessProbe(c echo.Context) error {
	ctx := context.Background()
	
	// Check database connection
	if err := h.db.Ping(ctx); err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"status": "not ready",
			"reason": "database connection failed",
		})
	}

	// Check Redis connection
	if err := h.cache.Client().Ping(ctx).Err(); err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"status": "not ready",
			"reason": "cache connection failed",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "ready",
	})
}
