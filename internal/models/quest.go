package models

import "time"

type Quest struct {
	ID             int    `json:"id" db:"id"`
	Title          string `json:"title" db:"title"`
	Description    string `json:"description" db:"description"`
	Category       string `json:"category" db:"category"`
	Rarity         string `json:"rarity" db:"rarity"`
	Difficulty     int    `json:"difficulty" db:"difficulty"`
	Price          int    `json:"price" db:"price"`
	RewardXP       int    `json:"reward_xp" db:"reward_xp"`
	RewardCoin     int    `json:"reward_coin" db:"reward_coin"`
	TimeLimitHours int    `json:"time_limit_hours" db:"time_limit_hours"`
	Tasks          []Task `json:"tasks,omitempty"`
}

type UserCurrentQuests struct {
	UserID      int       `json:"user_id" db:"user_id"`
	QuestID     int       `json:"quest_id" db:"quest_id"`
	Status      string    `json:"status" db:"status"` // "purchased", "started", "failed", "completed"
	TasksDone   int       `json:"tasks_done" db:"tasks_done"`
	StartedAt   time.Time `json:"started_at" db:"started_at"`
	CompletedAt time.Time `json:"completed_at" db:"completed_at"`
	ExpiresAt   time.Time `json:"expires_at" db:"expires_at"`
}

type Task struct {
	ID             int    `json:"id" db:"id"`
	Title          string `json:"title" db:"title"`
	Description    string `json:"description" db:"description"`
	Difficulty     int    `json:"difficulty" db:"difficulty"`
	Rarity         string `json:"rarity" db:"rarity"`
	Category       string `json:"category" db:"category"`
	BaseXpReward   int    `json:"base_xp_reward" db:"base_xp_reward"`
	BaseCoinReward int    `json:"base_coin_reward" db:"base_coin_reward"`
	// TODO: requeired level of health, intelligence, charisma, willpower ?
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
