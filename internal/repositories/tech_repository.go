package repositories

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type TechRepository struct {
	db *sqlx.DB
}

func NewTechRepository(db *sqlx.DB) *TechRepository {
	return &TechRepository{db: db}
}

func (r *TechRepository) CheckConnection() error {
	if err := r.db.Ping(); err != nil {
		return errors.New("failed to connect to the database")
	}
	return nil
}
