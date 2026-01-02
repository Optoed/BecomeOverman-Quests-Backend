package services

import "BecomeOverMan/internal/repositories"

type TechService struct {
	repo *repositories.TechRepository
}

func NewTechService(repo *repositories.TechRepository) *TechService {
	return &TechService{repo: repo}
}

func (s *TechService) CheckConnection() error {
	return s.repo.CheckConnection()
}
