package models

import (
	"encoding/json"
	"time"
)

type AIQuestResponse struct {
	Quest *Quest `json:"quest"`
	Tasks []Task `json:"tasks"`
}

type Quest struct {
	ID             int              `json:"id" db:"id"`
	Title          string           `json:"title" db:"title"`
	Description    string           `json:"description" db:"description"`
	Category       string           `json:"category" db:"category"`
	Rarity         string           `json:"rarity" db:"rarity"`
	Difficulty     int              `json:"difficulty" db:"difficulty"`
	Price          int              `json:"price" db:"price"`
	TasksCount     int              `json:"tasks_count" db:"tasks_count"`
	ConditionsJson *json.RawMessage `json:"conditions_json" db:"conditions_json"`
	BonusJson      *json.RawMessage `json:"bonus_json" db:"bonus_json"`
	IsSequential   bool             `json:"is_sequential" db:"is_sequential"`
	RewardXP       int              `json:"reward_xp" db:"reward_xp"`
	RewardCoin     int              `json:"reward_coin" db:"reward_coin"`
	TimeLimitHours int              `json:"time_limit_hours" db:"time_limit_hours"`
	Tasks          []Task           `json:"tasks,omitempty"`
}

type Task struct {
	ID             int       `json:"id" db:"id"`
	Title          string    `json:"title" db:"title"`
	Description    string    `json:"description" db:"description"`
	Difficulty     int       `json:"difficulty" db:"difficulty"`
	Rarity         string    `json:"rarity" db:"rarity"`
	Category       string    `json:"category" db:"category"`
	BaseXpReward   int       `json:"base_xp_reward" db:"base_xp_reward"`
	BaseCoinReward int       `json:"base_coin_reward" db:"base_coin_reward"`
	TaskOrder      int       `json:"task_order" db:"task_order"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`

	// --- опциональные поля для UserTask
	QuestID        *int       `json:"quest_id" db:"quest_id"`
	Status         *string    `json:"status" db:"status"` // nullable
	ScheduledStart *time.Time `json:"scheduled_start" db:"scheduled_start"`
	ScheduledEnd   *time.Time `json:"scheduled_end" db:"scheduled_end"`
	Deadline       *time.Time `json:"deadline" db:"deadline"`
	Duration       *int       `json:"duration" db:"duration"`
	UpdatedByAI    *bool      `json:"updated_by_ai" db:"updated_by_ai"`
	CompletedAt    *time.Time `json:"completed_at" db:"completed_at"`
	IsConfirmed    *bool      `json:"is_confirmed" db:"is_confirmed"`

	XpGained   *int `json:"xp_gained" db:"xp_gained"`     // nullable
	CoinGained *int `json:"coin_gained" db:"coin_gained"` // nullable
}

type UserQuests struct {
	UserID      int       `json:"user_id" db:"user_id"`
	QuestID     int       `json:"quest_id" db:"quest_id"`
	Status      string    `json:"status" db:"status"` // not_started, "active", "completed", TODO: failed
	StartedAt   time.Time `json:"started_at" db:"started_at"`
	CompletedAt time.Time `json:"completed_at" db:"completed_at"`
	ExpiresAt   time.Time `json:"expires_at" db:"expires_at"`

	TasksDone int `json:"tasks_done"`
}
