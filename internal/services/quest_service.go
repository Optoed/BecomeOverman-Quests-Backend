package services

import (
	"BecomeOverMan/internal/integrations"
	"BecomeOverMan/internal/models"
	"BecomeOverMan/internal/repositories"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type QuestService struct {
	questRepo *repositories.QuestRepository
}

func NewQuestService(repo *repositories.QuestRepository) *QuestService {
	return &QuestService{questRepo: repo}
}

// GetAvailableQuests returns quests available for the user
func (s *QuestService) GetAvailableQuests(ctx context.Context, userID int) ([]models.Quest, error) {
	return s.questRepo.GetAvailableQuests(ctx, userID)
}

func (s *QuestService) GetQuestShop(ctx context.Context, userID int) ([]models.Quest, error) {
	return s.questRepo.GetQuestShop(ctx, userID)
}

func (s *QuestService) GetMyActiveQuests(ctx context.Context, userID int) ([]models.Quest, error) {
	return s.questRepo.GetMyActiveQuests(ctx, userID)
}

func (s *QuestService) GetMyCompletedQuests(ctx context.Context, userID int) ([]models.Quest, error) {
	return s.questRepo.GetMyCompletedQuests(ctx, userID)
}

func (s *QuestService) GetMyAllQuestsWithDetails(ctx context.Context, userID int) ([]models.Quest, error) {
	return s.questRepo.GetMyAllQuestsWithDetails(ctx, userID)
}

// PurchaseQuest handles the purchase of a quest by a user
func (s *QuestService) PurchaseQuest(ctx context.Context, userID, questID int) error {
	return s.questRepo.PurchaseQuest(ctx, userID, questID)
}

// StartQuest begins the execution of a purchased quest
func (s *QuestService) StartQuest(ctx context.Context, userID, questID int) error {
	return s.questRepo.StartQuest(ctx, userID, questID)
}

// CompleteTask marks a task as completed by the user
func (s *QuestService) CompleteTask(ctx context.Context, userID, questID, taskID int) error {
	return s.questRepo.CompleteTask(ctx, userID, questID, taskID)
}

// CompleteQuest finalizes the quest completion
func (s *QuestService) CompleteQuest(ctx context.Context, userID, questID int) error {
	return s.questRepo.CompleteQuest(ctx, userID, questID)
}

func (s *QuestService) GetQuestDetails(ctx context.Context, questID int, userID int) (*models.Quest, error) {
	return s.questRepo.GetQuestDetails(ctx, questID, userID)
}

func (s *QuestService) CreateSharedQuest(user1ID, user2ID, questID int) error {
	return s.questRepo.CreateSharedQuest(user1ID, user2ID, questID)
}

func (s *QuestService) SaveQuestToDB(quest *models.Quest, tasks []models.Task) (int, error) {
	return s.questRepo.SaveQuestToDB(quest, tasks)
}

func (s *QuestService) SearchQuests(
	ctx context.Context,
	req models.RecommendationService_SearchQuest_Request,
) (models.SearchQuestsResponse, error) {
	// 1. Создаем URL
	url := integrations.Recommendation_Service_BASE_URL + "/search"

	// 2. Кодируем в JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// 3. Создаем io.Reader из JSON
	body := bytes.NewBuffer(jsonData)

	// 4. Делаем POST запрос
	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		return nil, fmt.Errorf("error making POST request to recommendation service: %v", err)
	}

	defer resp.Body.Close()

	// 5. Проверяем статус
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("recommendation service returned status %d", resp.StatusCode)
	}

	// 6. Читаем и парсим ответ
	var response models.RecommendationService_SearchQuests_Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// 7. достаем ids
	questsIDS := make([]int, len(response.Results))
	for i, result := range response.Results {
		questsIDS[i] = result.ID
	}

	// 8. Достаем квесты с деталями из БД (тут сразу и те что есть у юзера и те что еще не куплены)
	questsWithDetails, err := s.questRepo.SearchQuestsWithDetailsByIDs(ctx, questsIDS)
	if err != nil {
		slog.ErrorContext(ctx, "ошибка получения квестов из БД с указанными ids во время поиска",
			"error", err,
			"ids", questsIDS,
		)
		return nil, fmt.Errorf("В поиске квестов по запросу произошла внутренняя ошибка: %w", err)
	}

	// 9. Возвращаем результат = []struct{questWithDetails, SimilaryScore}
	questsWithDetailsAndSimilarityResponse := models.NewSearchQuestsResponse(questsWithDetails, response)

	return questsWithDetailsAndSimilarityResponse, nil
}
