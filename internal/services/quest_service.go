package services

import (
	"BecomeOverMan/internal/models"
	"BecomeOverMan/internal/repositories"
	"context"
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

// GenerateAIQuest вызывает Python скрипт для генерации квеста
func (s *QuestService) GenerateAIQuest(userMessage string) (*models.AIQuestResponse, error) {
	return GenerateAIQuest(userMessage)
}

func (s *QuestService) SaveQuestToDB(quest *models.Quest, tasks []models.Task) (int, error) {
	return s.questRepo.SaveQuestToDB(quest, tasks)
}
