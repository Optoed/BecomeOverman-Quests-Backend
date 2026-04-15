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

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=purchased active completed"`
}

func (h *QuestHandler) GetQuestDetails(c *gin.Context) {
	questIDStr := c.Param("questID")
	questID, err := strconv.Atoi(questIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quest ID"})
		return
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	questDetails, err := h.questService.GetQuestDetails(c.Request.Context(), questID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, questDetails)
}

// GetUserQuests returns user's quests, optionally filtered by ?status=active|completed
func (h *QuestHandler) GetUserQuests(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	status := c.Query("status")

	var quests []models.Quest
	switch status {
	case "active":
		quests, err = h.questService.GetMyActiveQuests(c.Request.Context(), userID)
	case "completed":
		quests, err = h.questService.GetMyCompletedQuests(c.Request.Context(), userID)
	default:
		quests, err = h.questService.GetMyAllQuestsWithDetails(c.Request.Context(), userID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quests)
}

func (h *QuestHandler) GetAvailableQuestsHandler(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	quests, err := h.questService.GetAvailableQuests(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quests)
}

func (h *QuestHandler) GetQuestShopHandler(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	quests, err := h.questService.GetQuestShop(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quests)
}

// UpdateQuestStatus handles PATCH /users/me/quests/:questID — purchase, start, or complete a quest
func (h *QuestHandler) UpdateQuestStatus(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	questID, err := strconv.Atoi(c.Param("questID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quest ID"})
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: status must be one of: purchased, active, completed"})
		return
	}

	switch req.Status {
	case "purchased":
		err = h.questService.PurchaseQuest(c.Request.Context(), userID, questID)
	case "active":
		err = h.questService.StartQuest(c.Request.Context(), userID, questID)
	case "completed":
		err = h.questService.CompleteQuest(c.Request.Context(), userID, questID)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"quest_id": questID, "status": req.Status})
}

// UpdateTaskStatus handles PATCH /users/me/quests/:questID/tasks/:taskID — complete a task
func (h *QuestHandler) UpdateTaskStatus(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
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

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Status != "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status for task. Allowed: completed"})
		return
	}

	if err := h.questService.CompleteTask(c.Request.Context(), userID, questID, taskID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"quest_id": questID, "task_id": taskID, "status": req.Status})
}

func (h *QuestHandler) CreateSharedQuest(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
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

	c.JSON(http.StatusCreated, gin.H{"message": "Shared quest created successfully"})
}

// SearchQuests handles GET /quests/search?q=...&top_k=...&category=...&status=...
func (h *QuestHandler) SearchQuests(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	topK := 5
	if topKStr := c.Query("top_k"); topKStr != "" {
		if val, err := strconv.Atoi(topKStr); err == nil && val > 0 && val <= 100 {
			topK = val
		}
	}

	req := models.RecommendationService_SearchQuest_Request{
		Query:    query,
		TopK:     topK,
		Category: c.Query("category"),
		Status:   c.DefaultQuery("status", "all"),
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	quests, err := h.questService.SearchQuests(c.Request.Context(), req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quests)
}

func (h *QuestHandler) RecommendFriends(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	req := models.RecommendationService_RecommendUsers_Request{
		UserID: userID,
	}

	resp, err := h.questService.RecommendFriends(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *QuestHandler) RecommendQuests(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.questService.RecommendQuests(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func RegisterQuestRoutes(router *gin.Engine, questService *services.QuestService) {
	handler := NewQuestHandler(questService)

	questGroup := router.Group("/quests")
	questGroup.Use(middleware.JWTAuthMiddleware())
	{
		questGroup.GET("/available", handler.GetAvailableQuestsHandler)
		questGroup.GET("/shop", handler.GetQuestShopHandler)
		questGroup.GET("/search", handler.SearchQuests)
		questGroup.GET("/:questID", handler.GetQuestDetails)

		questGroup.POST("", handler.GenerateAIQuest)
		questGroup.POST("/shared", handler.CreateSharedQuest)
	}

	userQuestsGroup := router.Group("/users/me")
	userQuestsGroup.Use(middleware.JWTAuthMiddleware())
	{
		userQuestsGroup.GET("/quests", handler.GetUserQuests)
		userQuestsGroup.PATCH("/quests/:questID", handler.UpdateQuestStatus)
		userQuestsGroup.PATCH("/quests/:questID/tasks/:taskID", handler.UpdateTaskStatus)
		userQuestsGroup.GET("/recommendations/quests", handler.RecommendQuests)
		userQuestsGroup.GET("/recommendations/friends", handler.RecommendFriends)
	}

	scheduleGroup := router.Group("/schedules")
	scheduleGroup.Use(middleware.JWTAuthMiddleware())
	{
		scheduleGroup.POST("", handler.GenerateScheduleByAI)
	}
}
