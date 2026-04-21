package repositories

import (
	"BecomeOverMan/internal/models"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

func (r *UserRepository) CreateUserWithProfile(username, email, hashedPassword string) (models.UserProfile, error) {
	var user models.UserProfile
	query := `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, username, email, xp_points, coin_balance, level, current_streak, longest_streak, created_at, last_active_at, version
	`
	err := r.db.Get(&user, query, username, email, hashedPassword)
	if err != nil {
		return models.UserProfile{}, err
	}
	return user, nil
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
	query := `SELECT id, username, email, version, xp_points, coin_balance, level, created_at FROM users WHERE id = $1`
	err := r.db.Get(&user, query, userID)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByID(userID int) (models.UserProfile, error) {
	var user models.UserProfile
	query := `
		SELECT id, username, email, xp_points, coin_balance, level, current_streak, longest_streak, created_at, last_active_at, version
		FROM users
		WHERE id = $1
	`
	err := r.db.Get(&user, query, userID)
	if err != nil {
		return models.UserProfile{}, err
	}
	return user, nil
}

func (r *UserRepository) GetProfiles(userIDs []int) ([]models.UserProfile, error) {
	if len(userIDs) == 0 {
		return []models.UserProfile{}, nil
	}

	var usersProfiles []models.UserProfile
	query := `
		SELECT id, username, email, version, xp_points, coin_balance, level, current_streak, longest_streak, created_at, last_active_at
		FROM users
		WHERE id = ANY($1)
		ORDER BY array_position($1, id)
		`
	err := r.db.Select(&usersProfiles, query, pq.Array(userIDs))
	if err != nil {
		return nil, err
	}

	return usersProfiles, nil
}

func (r *UserRepository) ListUsers(limit, offset int) ([]models.UserProfile, error) {
	var users []models.UserProfile
	query := `
		SELECT id, username, email, version, xp_points, coin_balance, level, current_streak, longest_streak, created_at, last_active_at
		FROM users
		ORDER BY id
		LIMIT $1 OFFSET $2
	`
	if err := r.db.Select(&users, query, limit, offset); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) UpdateUser(id int, req models.UpdateUserRequest) (models.UserProfile, error) {
	var updated models.UserProfile
	query := `
		UPDATE users
		SET
			username = COALESCE($1, username),
			email = COALESCE($2, email),
			version = version + 1,
			last_active_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND version = $4
		RETURNING id, username, email, version, xp_points, coin_balance, level, current_streak, longest_streak, created_at, last_active_at
	`
	err := r.db.Get(&updated, query, req.Username, req.Email, id, req.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.UserProfile{}, sql.ErrNoRows
		}
		return models.UserProfile{}, err
	}
	return updated, nil
}

func (r *UserRepository) DeleteUser(id int) (bool, error) {
	res, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}
