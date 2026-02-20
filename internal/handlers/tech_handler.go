package handlers

import (
	"BecomeOverMan/internal/integrations"
	"BecomeOverMan/internal/services"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TechHandler struct {
	service *services.TechService
}

func NewTechHandler(service *services.TechService) *TechHandler {
	return &TechHandler{service: service}
}

func (h *TechHandler) CheckConnectionDB(c *gin.Context) {
	if err := h.service.CheckConnection(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "database is available"})
}

func (h *TechHandler) RecommendationServiceHealth(c *gin.Context) {
	url := integrations.Recommendation_Service_BASE_URL + "/health"

	resp, err := http.Get(url)
	if err != nil {
		slog.Error("Failed to make request to recommendation service", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Cannot connect to recommendation service",
			"details": err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body: " + err.Error()})
		return
	}

	slog.Debug("Response body", "body", string(body))

	// Парсим JSON ответ
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(resp.StatusCode, result)
}

func RegisterTechRoutes(router *gin.Engine, techService *services.TechService) {
	handler := NewTechHandler(techService)

	g := router.Group("/tech")
	{
		g.GET("/ping-db", handler.CheckConnectionDB)

		g.GET("/recommendation-service/health", handler.RecommendationServiceHealth)
	}
}
