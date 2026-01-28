package repositories

import (
	"BecomeOverMan/internal/models"
	"errors"

	"github.com/jmoiron/sqlx"
)

func (r *UserRepository) AddFriend(userID, friendID int) error {
	// Проверяем, что пользователь существует
	var userExists bool
	err := r.db.Get(&userExists, `
		SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`, friendID)
	if err != nil {
		return err
	}
	if !userExists {
		return errors.New("user not found")
	}

	// Проверяем, что дружба не существует
	var friendshipExists bool
	err = r.db.Get(&friendshipExists, `
		SELECT EXISTS(SELECT 1 FROM friends WHERE user_id = $1 AND friend_id = $2)`,
		userID, friendID)
	if err != nil {
		return err
	}
	if friendshipExists {
		return errors.New("friendship already exists")
	}

	_, err = r.db.Exec(`
		INSERT INTO friends (user_id, friend_id, status) 
		VALUES ($1, $2, 'accepted')`,
		userID, friendID)
	return err
}

func (r *UserRepository) GetFriends(userID int) ([]models.Friend, error) {
	var friends []models.Friend
	query := `
		SELECT f.*, u.username 
		FROM friends f 
		JOIN users u ON f.friend_id = u.id 
		WHERE f.user_id = $1 AND f.status = 'accepted'
	`
	err := r.db.Select(&friends, query, userID)
	return friends, err
}

// Shared quest methods
func (r *QuestRepository) CreateSharedQuest(user1ID, user2ID, questID int) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Проверяем, что пользователи друзья
	var areFriends bool
	err = tx.Get(&areFriends, `
		SELECT EXISTS(
			SELECT 1 FROM friends 
			WHERE user_id = $1 AND friend_id = $2 AND status = 'accepted'
		)`, user1ID, user2ID)
	if err != nil {
		return err
	}
	if !areFriends {
		return errors.New("users are not friends")
	}

	// Создаем shared quest
	_, err = tx.Exec(`
		INSERT INTO shared_quests (user1_id, user2_id, quest_id, status) 
		VALUES ($1, $2, $3, 'active')`,
		user1ID, user2ID, questID)
	if err != nil {
		return err
	}

	// Стартуем квест для обоих пользователей
	if err := r.startQuestForUser(tx, user1ID, questID); err != nil {
		return err
	}
	if err := r.startQuestForUser(tx, user2ID, questID); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *QuestRepository) startQuestForUser(tx *sqlx.Tx, userID, questID int) error {
	// Покупаем квест (если еще не куплен)
	var alreadyPurchased bool
	err := tx.Get(&alreadyPurchased, `
		SELECT EXISTS(SELECT 1 FROM user_quests WHERE user_id = $1 AND quest_id = $2)`,
		userID, questID)
	if err != nil {
		return err
	}

	if alreadyPurchased {
		return errors.New("quest already purchased")
	}

	// Получаем цену квеста
	var price int
	err = tx.Get(&price, "SELECT price FROM quests WHERE id = $1", questID)
	if err != nil {
		return err
	}

	// Проверяем баланс
	var balance int
	err = tx.Get(&balance, "SELECT coin_balance FROM users WHERE id = $1", userID)
	if err != nil {
		return err
	}

	if balance < price {
		return errors.New("not enough coins for shared quest")
	}

	// Покупаем квест
	_, err = tx.Exec(`
			INSERT INTO user_quests (user_id, quest_id, status) 
			VALUES ($1, $2, 'purchased')`,
		userID, questID)
	if err != nil {
		return err
	}

	// Списываем монеты
	_, err = tx.Exec(`
			UPDATE users SET coin_balance = coin_balance - $1 WHERE id = $2`,
		price, userID)
	if err != nil {
		return err
	}

	// Стартуем квест
	_, err = tx.Exec(`
			UPDATE user_quests 
			SET status = 'started', started_at = NOW() 
			WHERE user_id = $1 AND quest_id = $2`,
		userID, questID)

	return err
}
