package handlers

import (
	"BecomeOverMan/pkg/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestAI struct {
	Prompt string `json:"prompt" binding:"required"`
}

func (h *QuestHandler) GenerateAIQuest(c *gin.Context) {
	// тут из запроса пользователя достаем текст что он написал во фротенде для генерации ему квеста
	var request RequestAI
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// тут вызываем из питона эту функцию с ai и одновременно сохраняем в БД пользователю и одновременно возвращаем на фронтенд
	aiResponse, err := h.questService.GenerateAIQuest(request.Prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate quest: " + err.Error()})
		return
	}

	// Сохраняем квест в БД
	questID, err := h.questService.SaveQuestToDB(aiResponse.Quest, aiResponse.Tasks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save quest: " + err.Error()})
		return
	}

	// Возвращаем ответ на фронтенд
	c.JSON(http.StatusOK, gin.H{
		"message":  "Quest generated successfully",
		"quest_id": questID,
		"quest":    aiResponse.Quest,
		"tasks":    aiResponse.Tasks,
	})
}

func (h *QuestHandler) GenerateScheduleByAI(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var request RequestAI
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// в том числе сохранили в БД
	aiResponse, err := h.questService.GenerateScheduleByAI(c.Request.Context(), userID, request.Prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate schedule: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, aiResponse)
}
