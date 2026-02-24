package handlers

import (
	"BecomeOverMan/internal/integrations"
	"BecomeOverMan/internal/models"
	"BecomeOverMan/pkg/middleware"
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

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

	// параллельно отправляем в Сервис Рекоммендаций
	req := models.RecommendationService_AddQuests_Request{
		Quests: []models.RecommendationService_questToAdd{
			{
				ID:          strconv.Itoa(questID),
				Title:       aiResponse.Quest.Title,
				Description: aiResponse.Quest.Description,
				Category:    aiResponse.Quest.Category,
			},
		},
	}

	go func() {
		err := h.sendQuestToRecommendationService(req)
		if err != nil {
			slog.Error("Failed to send (add) quest to recommendation service", "error", err)
		}
	}()

	// Возвращаем ответ на фронтенд
	c.JSON(http.StatusOK, gin.H{
		"message":  "Quest generated successfully",
		"quest_id": questID,
		"quest":    aiResponse.Quest,
		"tasks":    aiResponse.Tasks,
	})
}

func (h *QuestHandler) sendQuestToRecommendationService(req models.RecommendationService_AddQuests_Request) error {
	// 1. Создаем URL
	url := integrations.Recommendation_Service_BASE_URL + "/quests/add"

	// 2. Кодируем в JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	// 3. Создаем io.Reader из JSON
	body := bytes.NewBuffer(jsonData)

	// 4. Делаем POST запрос
	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		return fmt.Errorf("error making POST request (add quest) to recommendation service: %v", err)
	}

	defer resp.Body.Close()

	// 5. Проверяем статус
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("recommendation service returned status %d", resp.StatusCode)
	}

	return nil
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
