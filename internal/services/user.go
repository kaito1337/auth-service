package services

import (
	"auth-backend/internal/models"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	pgClient *gorm.DB
}

func NewUserService(pgClient *gorm.DB) *UserService {
	return &UserService{
		pgClient: pgClient,
	}
}

func (s *UserService) GetUserById(id string) (*models.User, error) {
	var user *models.User
	result := s.pgClient.Model(&models.User{}).Where("id = ?", id).First(&user)
	if result.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user *models.User
	result := s.pgClient.Model(&models.User{}).Where("email = ?", email).First(&user)
	if result.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *UserService) GetUserByLogin(login string) (*models.User, error) {
	var user *models.User
	result := s.pgClient.Model(&models.User{}).Where("login = ?", login).First(&user)
	if result.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *UserService) CreateUser(email string, login string, password string) (*models.User, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	var user = &models.User{
		Login:    login,
		Email:    email,
		Password: string(hashedPassword),
	}
	if _, err := s.GetUserByEmail(email); err == nil {
		return nil, errors.New("user with this email already exists")
	}
	if _, err := s.GetUserByLogin(login); err == nil {
		return nil, errors.New("user with this login already exists")
	}
	result := s.pgClient.Create(user)
	if result.RowsAffected == 0 {
		return nil, result.Error
	}
	return user, nil
}

func (s *UserService) CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *UserService) UpdateRefreshToken(email string, refreshToken string) error {
	result := s.pgClient.Model(&models.User{}).Where("email = ?", email).Update("refresh_token", refreshToken)
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}
