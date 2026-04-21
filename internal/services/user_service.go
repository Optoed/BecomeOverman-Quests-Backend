package services

import (
	"BecomeOverMan/internal/models"
	"BecomeOverMan/internal/repositories"
	"database/sql"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repositories.UserRepository
}

var ErrUserVersionConflict = errors.New("user version conflict")

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(username, email, password string) error {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	hashedPassword := string(hashedPasswordBytes)
	return s.repo.CreateUser(username, email, hashedPassword)
}

func (s *UserService) CreateUser(req models.CreateUserRequest) (models.UserProfile, error) {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.UserProfile{}, err
	}
	return s.repo.CreateUserWithProfile(req.Username, req.Email, string(hashedPasswordBytes))
}

func (s *UserService) Login(username, password string) (int, error) {
	// Логируем начало попытки входа
	log.Printf("User %s is attempting to log in", username)

	// Получаем пользователя по username
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		// Логируем ошибку при получении пользователя
		log.Printf("Error getting user by username %s: %v", username, err)
		return 0, errors.New("invalid username or password")
	}

	// Сравниваем пароли
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		// Логируем неудачную попытку входа
		log.Printf("Failed login attempt for user %s: invalid password", username)
		return 0, errors.New("invalid username or password")
	}

	// Логируем успешный вход
	log.Printf("User %s logged in successfully", username)

	return user.ID, nil
}

func (s *UserService) AddFriend(userID, friendID int) error {
	return s.repo.AddFriend(userID, friendID)
}

func (s *UserService) AddFriendByName(userID int, friendName string) error {
	return s.repo.AddFriendbyName(userID, friendName)
}

func (s *UserService) GetFriends(userID int) ([]models.Friend, error) {
	return s.repo.GetFriends(userID)
}

func (s *UserService) GetProfile(userID int) (models.User, error) {
	return s.repo.GetProfile(userID)
}

func (s *UserService) GetUserByID(userID int) (models.UserProfile, error) {
	return s.repo.GetUserByID(userID)
}

func (s *UserService) ListUsers(limit, offset int) ([]models.UserProfile, error) {
	return s.repo.ListUsers(limit, offset)
}

func (s *UserService) UpdateUser(userID int, req models.UpdateUserRequest) (models.UserProfile, error) {
	updated, err := s.repo.UpdateUser(userID, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, existsErr := s.repo.GetUserByID(userID)
			if existsErr != nil {
				return models.UserProfile{}, existsErr
			}
			return models.UserProfile{}, ErrUserVersionConflict
		}
		return models.UserProfile{}, err
	}
	return updated, nil
}

func (s *UserService) DeleteUser(userID int) (bool, error) {
	return s.repo.DeleteUser(userID)
}
