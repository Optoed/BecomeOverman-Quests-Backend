package repositories

import (
	"context"
	"errors"
	"time"

	"BecomeOverMan/internal/models"

	"github.com/jmoiron/sqlx"
)

type QuestRepository struct {
	db *sqlx.DB
}

func NewQuestRepository(db *sqlx.DB) *QuestRepository {
	return &QuestRepository{db: db}
}

// GetAvailableQuests возвращает квесты, доступные для пользователя
func (r *QuestRepository) GetAvailableQuests(ctx context.Context, userID int) ([]models.Quest, error) {
	var quests []models.Quest

	// Получаем уровень пользователя
	var user models.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", userID)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT q.* FROM quests q 
		WHERE q.difficulty <= $1 + 1 AND q.price <= $2
	`

	err = r.db.SelectContext(ctx, &quests, query, user.Level, user.CoinBalance)
	if err != nil {
		return nil, err
	}

	return quests, nil
}

// PurchaseQuest покупает квест для пользователя
func (r *QuestRepository) PurchaseQuest(ctx context.Context, userID, questID int) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Проверяем, что квест доступен
	var quest models.Quest
	err = tx.GetContext(ctx, &quest, "SELECT * FROM quests WHERE id = $1", questID)
	if err != nil {
		return err
	}

	// Проверяем баланс пользователя
	var balance int
	err = tx.GetContext(ctx, &balance, "SELECT coin_balance FROM users WHERE id = $1", userID)
	if err != nil {
		return err
	}

	if balance < quest.Price {
		return errors.New("not enough currency")
	}

	// Списываем валюту
	_, err = tx.ExecContext(ctx,
		"UPDATE users SET coin_balance = coin_balance - $1 WHERE id = $2",
		quest.Price, userID)
	if err != nil {
		return err
	}

	// Добавляем квест пользователю
	_, err = tx.ExecContext(ctx, `
        INSERT INTO user_current_quests 
        (user_id, quest_id, status, started_at, expires_at)
        VALUES ($1, $2, 'purchased', NULL, NULL)`,
		userID, questID)
	if err != nil {
		return err
	}

	// Записываем транзакцию
	_, err = tx.ExecContext(ctx, `
        INSERT INTO user_coin_transactions 
        (user_id, amount, transaction_type, reference_type, reference_id, description)
        VALUES ($1, $2, 'spent', 'quest', $3, 'Purchased quest: ' || $4)`,
		userID, -quest.Price, quest.ID, quest.Title)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// StartQuest начинает выполнение квеста
func (r *QuestRepository) StartQuest(ctx context.Context, userID, questID int) error {
	// Устанавливаем время начала и завершения
	var timeLimitHours int
	err := r.db.GetContext(ctx, &timeLimitHours,
		"SELECT time_limit_hours FROM quests WHERE id = $1", questID)
	if err != nil {
		return err
	}

	expiresAt := time.Now().Add(time.Duration(timeLimitHours) * time.Hour)

	_, err = r.db.ExecContext(ctx, `
        UPDATE user_current_quests 
        SET status = 'started', started_at = NOW(), expires_at = $1
        WHERE user_id = $2 AND quest_id = $3 AND status = 'purchased'`,
		expiresAt, userID, questID)

	return err
}

// CompleteTask отмечает выполнение задачи
func (r *QuestRepository) CompleteTask(ctx context.Context, userID, taskID int) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Проверяем, что задача принадлежит активному квесту
	var questID int
	err = tx.GetContext(ctx, &questID, `
        SELECT q.id FROM quests q
        JOIN quest_tasks qt ON q.id = qt.quest_id
        JOIN user_current_quests ucq ON q.id = ucq.quest_id
        WHERE ucq.user_id = $1 AND qt.task_id = $2 AND ucq.status = 'started'`,
		userID, taskID)
	if err != nil {
		return err
	}

	// Добавляем задачу в буфер выполненных
	_, err = tx.ExecContext(ctx, `
        INSERT INTO user_completed_tasks 
        (user_id, task_id, completed_at, xp_gained, coin_gained)
        SELECT $1, $2, NOW(), t.base_xp_reward, t.base_coin_reward
        FROM tasks t WHERE t.id = $2`,
		userID, taskID)
	if err != nil {
		return err
	}

	// Увеличиваем счетчик выполненных задач
	_, err = tx.ExecContext(ctx, `
        UPDATE user_current_quests 
        SET tasks_done = tasks_done + 1
        WHERE user_id = $1 AND quest_id = $2`,
		userID, questID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// CompleteQuest завершает квест
func (r *QuestRepository) CompleteQuest(ctx context.Context, userID, questID int) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Проверяем, что все задачи выполнены
	var totalTasks, completedTasks int
	err = tx.GetContext(ctx, &totalTasks, `
        SELECT COUNT(*) FROM quest_tasks WHERE quest_id = $1`, questID)
	if err != nil {
		return err
	}

	err = tx.GetContext(ctx, &completedTasks, `
        SELECT tasks_done FROM user_current_quests 
        WHERE user_id = $1 AND quest_id = $2`,
		userID, questID)
	if err != nil {
		return err
	}

	if completedTasks < totalTasks {
		return errors.New("not all tasks completed")
	}

	// Получаем награду за квест
	var rewardXP, rewardCoin int
	err = tx.QueryRowContext(ctx, `
    	SELECT reward_xp, reward_coin FROM quests WHERE id = $1`, questID).
		Scan(&rewardXP, &rewardCoin)
	if err != nil {
		return err
	}

	// Начисляем награды
	_, err = tx.ExecContext(ctx, `
        UPDATE users 
        SET xp_points = xp_points + $1,
            coin_balance = coin_balance + $2
        WHERE id = $3`,
		rewardXP, rewardCoin, userID)
	if err != nil {
		return err
	}

	// Переносим квест в завершенные
	_, err = tx.ExecContext(ctx, `
        INSERT INTO user_completed_quests 
        (user_id, quest_id, completed_at, xp_gained, coin_gained)
        VALUES ($1, $2, NOW(), $3, $4)`,
		userID, questID, rewardXP, rewardCoin)
	if err != nil {
		return err
	}

	// Удаляем из активных
	_, err = tx.ExecContext(ctx, `
        DELETE FROM user_current_quests 
        WHERE user_id = $1 AND quest_id = $2`,
		userID, questID)
	if err != nil {
		return err
	}

	// Переносим задачи из буфера в историю
	_, err = tx.ExecContext(ctx, `
        UPDATE user_completed_tasks 
        SET is_confirmed = true
        WHERE user_id = $1 AND task_id IN (
            SELECT task_id FROM quest_tasks WHERE quest_id = $2
        )`,
		userID, questID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
