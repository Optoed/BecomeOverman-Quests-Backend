package handlers

import (
	"BecomeOverMan/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QuestHandler struct {
	questService *services.QuestService
}

func NewQuestHandler(questService *services.QuestService) *QuestHandler {
	return &QuestHandler{questService: questService}
}

// GetAvailableQuestsHandler handles the GET request for available quests
// @Summary Get available quests for user
// @Description Returns quests available for the current user based on their level and coin balance
// @Tags Quests
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header int true "User ID"
// @Success 200 {array} models.Quest
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /quests/available [get]
func (h *QuestHandler) GetAvailableQuestsHandler(c *gin.Context) {
	userID, err := strconv.Atoi(c.GetHeader("X-User-ID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	quests, err := h.questService.GetAvailableQuests(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quests)
}

// PurchaseQuestHandler handles the POST request to purchase a quest
// @Summary Purchase a quest
// @Description Allows user to purchase a quest if they have enough currency
// @Tags Quests
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header int true "User ID"
// @Param questID path int true "Quest ID"
// @Success 200
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /quests/{questID}/purchase [post]
func (h *QuestHandler) PurchaseQuestHandler(c *gin.Context) {
	userID, err := strconv.Atoi(c.GetHeader("X-User-ID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	questID, err := strconv.Atoi(c.Param("questID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quest ID"})
		return
	}

	if err := h.questService.PurchaseQuest(c.Request.Context(), userID, questID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// StartQuestHandler handles the POST request to start a quest
// @Summary Start a purchased quest
// @Description Begins the execution of a purchased quest
// @Tags Quests
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header int true "User ID"
// @Param questID path int true "Quest ID"
// @Success 200
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /quests/{questID}/start [post]
func (h *QuestHandler) StartQuestHandler(c *gin.Context) {
	userID, err := strconv.Atoi(c.GetHeader("X-User-ID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	questID, err := strconv.Atoi(c.Param("questID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quest ID"})
		return
	}

	if err := h.questService.StartQuest(c.Request.Context(), userID, questID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// CompleteTaskHandler handles the POST request to complete a task
// @Summary Complete a quest task
// @Description Marks a specific task as completed by the user
// @Tags Quests
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header int true "User ID"
// @Param taskID path int true "Task ID"
// @Success 200
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{taskID}/complete [post]
func (h *QuestHandler) CompleteTaskHandler(c *gin.Context) {
	userID, err := strconv.Atoi(c.GetHeader("X-User-ID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	taskID, err := strconv.Atoi(c.Param("taskID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	if err := h.questService.CompleteTask(c.Request.Context(), userID, taskID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// CompleteQuestHandler handles the POST request to complete a quest
// @Summary Complete a quest
// @Description Finalizes quest completion if all tasks are done
// @Tags Quests
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header int true "User ID"
// @Param questID path int true "Quest ID"
// @Success 200
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /quests/{questID}/complete [post]
func (h *QuestHandler) CompleteQuestHandler(c *gin.Context) {
	userID, err := strconv.Atoi(c.GetHeader("X-User-ID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	questID, err := strconv.Atoi(c.Param("questID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quest ID"})
		return
	}

	if err := h.questService.CompleteQuest(c.Request.Context(), userID, questID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// RegisterQuestRoutes sets up the routes for quest handling with Gin
func RegisterQuestRoutes(router *gin.RouterGroup, questService *services.QuestService) {
	handler := NewQuestHandler(questService)

	questGroup := router.Group("/quests")
	{
		questGroup.GET("/available", handler.GetAvailableQuestsHandler)
		questGroup.POST("/:questID/purchase", handler.PurchaseQuestHandler)
		questGroup.POST("/:questID/start", handler.StartQuestHandler)
		questGroup.POST("/:questID/complete", handler.CompleteQuestHandler)
	}

	taskGroup := router.Group("/tasks")
	{
		taskGroup.POST("/:taskID/complete", handler.CompleteTaskHandler)
	}
}
