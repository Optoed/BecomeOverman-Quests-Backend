package repositories

import (
	"BecomeOverMan/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(username, email, hashedPassword string) error {
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, username, email, hashedPassword)
	return err
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, password_hash FROM users WHERE username = $1`
	err := r.db.Get(&user, query, username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetProfile(userID int) (models.User, error) {
	var user models.User
	query := `SELECT id, username, email, xp_points, coin_balance, level, created_at FROM users WHERE id = $1`
	err := r.db.Get(&user, query, userID)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
