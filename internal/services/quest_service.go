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

// PurchaseQuest handles the purchase of a quest by a user
func (s *QuestService) PurchaseQuest(ctx context.Context, userID, questID int) error {
	return s.questRepo.PurchaseQuest(ctx, userID, questID)
}

// StartQuest begins the execution of a purchased quest
func (s *QuestService) StartQuest(ctx context.Context, userID, questID int) error {
	return s.questRepo.StartQuest(ctx, userID, questID)
}

// CompleteTask marks a task as completed by the user
func (s *QuestService) CompleteTask(ctx context.Context, userID, taskID int) error {
	return s.questRepo.CompleteTask(ctx, userID, taskID)
}

// CompleteQuest finalizes the quest completion
func (s *QuestService) CompleteQuest(ctx context.Context, userID, questID int) error {
	return s.questRepo.CompleteQuest(ctx, userID, questID)
}
