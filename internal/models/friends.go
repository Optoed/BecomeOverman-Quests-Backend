package models

import "time"

type Friend struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	FriendID  int       `json:"friend_id" db:"friend_id"`
	Status    string    `json:"status" db:"status"`
	Username  string    `json:"username" db:"username"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type SharedQuest struct {
	ID        int       `json:"id" db:"id"`
	QuestID   int       `json:"quest_id" db:"quest_id"`
	User1ID   int       `json:"user1_id" db:"user1_id"`
	User2ID   int       `json:"user2_id" db:"user2_id"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CreateSharedQuestRequest struct {
	FriendID int `json:"friend_id" binding:"required"`
	QuestID  int `json:"quest_id" binding:"required"`
}

type AddFriendRequest struct {
	FriendID   *int    `json:"friend_id"`
	FriendName *string `json:"friend_name"`
}
