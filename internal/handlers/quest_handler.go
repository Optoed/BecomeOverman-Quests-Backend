package handlers

import (
	"BecomeOverMan/internal/models"
	"BecomeOverMan/internal/services"
	"BecomeOverMan/pkg/middleware"
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

// ==== Handlers ====

// GetQuestDetails возвращает детали квеста с задачами
func (h *QuestHandler) GetQuestDetails(c *gin.Context) {
	questIDStr := c.Param("questID")
	questID, err := strconv.Atoi(questIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quest ID"})
		return
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	questDetails, err := h.questService.GetQuestDetails(c.Request.Context(), questID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quest not found"})
		return
	}

	c.JSON(http.StatusOK, questDetails)
}

// GetAvailableQuestsHandler handles the GET request for available quests
// @Summary Get available quests for user
// @Description Returns quests available for the current user based on their level and coin balance
// @Tags Quests
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} models.Quest
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /quests/available [get]
func (h *QuestHandler) GetAvailableQuestsHandler(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	quests, err := h.questService.GetAvailableQuests(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quests)
}

// GetQuestShopHandler handles the GET request for quest shop
// @Summary Get quest shop for user
// @Description Returns quest shop
// @Tags Quests
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} models.Quest
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /quests/shop [get]
func (h *QuestHandler) GetQuestShopHandler(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	quests, err := h.questService.GetQuestShop(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quests)
}

// GetMyActiveQuests handles the GET request for active quests
// @Summary Get active quests for user
// @Description Returns active quests for user
// @Tags Quests
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} models.Quest
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /quests/active [get]
func (h *QuestHandler) GetMyActiveQuestsHandler(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	quests, err := h.questService.GetMyActiveQuests(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quests)
}

// GetMyCompletedQuests handles the GET request for completed quests
// @Summary Get completed quests for user
// @Description Returns completed quests for user
// @Tags Quests
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} models.Quest
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /quests/completed [get]
func (h *QuestHandler) GetMyCompletedQuests(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	quests, err := h.questService.GetMyCompletedQuests(c.Request.Context(), userID)
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
// @Param questID path int true "Quest ID"
// @Success 200
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /quests/{questID}/purchase [post]
func (h *QuestHandler) PurchaseQuestHandler(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
// @Param questID path int true "Quest ID"
// @Success 200
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /quests/{questID}/start [post]
func (h *QuestHandler) StartQuestHandler(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
// @Param questID path int true "Quest ID"
// @Param taskID path int true "Task ID"
// @Success 200
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /quests/{questID}/{taskID}/complete [post]
func (h *QuestHandler) CompleteTaskHandler(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	questID, err := strconv.Atoi(c.Param("questID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quest ID"})
		return
	}

	taskID, err := strconv.Atoi(c.Param("taskID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	if err := h.questService.CompleteTask(c.Request.Context(), userID, questID, taskID); err != nil {
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
// @Param questID path int true "Quest ID"
// @Success 200
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /quests/{questID}/complete [post]
func (h *QuestHandler) CompleteQuestHandler(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

// CreateSharedQuest godoc
// @Summary Create shared quest
// @Description Create a shared quest with a friend
// @Tags quests
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body models.CreateSharedQuestRequest true "Shared quest request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /quests/shared [post]
func (h *QuestHandler) CreateSharedQuest(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req models.CreateSharedQuestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.questService.CreateSharedQuest(userID, req.FriendID, req.QuestID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shared quest created successfully"})
}

// RegisterQuestRoutes sets up the routes for quest handling with Gin
func RegisterQuestRoutes(router *gin.Engine, questService *services.QuestService) {
	handler := NewQuestHandler(questService)

	questGroup := router.Group("/quests")
	questGroup.Use(middleware.JWTAuthMiddleware())
	{
		questGroup.GET("/:questID/details", handler.GetQuestDetails)
		questGroup.GET("/available", handler.GetAvailableQuestsHandler)
		questGroup.GET("/shop", handler.GetQuestShopHandler)
		questGroup.GET("/active", handler.GetMyActiveQuestsHandler)
		questGroup.GET("/completed", handler.GetMyCompletedQuests)
		questGroup.POST("/:questID/purchase", handler.PurchaseQuestHandler)
		questGroup.POST("/:questID/start", handler.StartQuestHandler)
		questGroup.POST("/:questID/complete", handler.CompleteQuestHandler)
		questGroup.POST("/:questID/:taskID/complete", handler.CompleteTaskHandler)

		questGroup.POST("/shared", handler.CreateSharedQuest)
	}
}
