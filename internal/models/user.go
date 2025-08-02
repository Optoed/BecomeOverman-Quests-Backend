package models

import "time"

type User struct {
	ID                int       `json:"id" db:"id"`
	Username          string    `json:"username" db:"username"`
	Email             string    `json:"email" db:"email"`
	PasswordHash      string    `json:"password_hash" db:"password_hash"`
	IsPremium         bool      `json:"is_premium" db:"is_premium"`
	XpPoints          int       `json:"xp_points" db:"xp_points"`
	CoinBalance       int       `json:"coin_balance" db:"coin_balance"`
	HealthLevel       int       `json:"health_level" db:"health_level"`
	IntelligenceLevel int       `json:"intelligence_level" db:"intelligence_level"`
	CharismaLevel     int       `json:"charisma_level" db:"charisma_level"`
	WillpowerLevel    int       `json:"willpower_level" db:"willpower_level"`
	Level             int       `json:"level" db:"level"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	LastActiveAt      time.Time `json:"last_active_at" db:"last_active_at"`
	CurrentStreak     int       `json:"current_streak" db:"current_streak"`
	LongestStreak     int       `json:"longest_streak" db:"longest_streak"`
}
