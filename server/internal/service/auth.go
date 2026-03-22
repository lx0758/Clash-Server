package service

import (
	"errors"

	"clash-server/internal/model"
	"clash-server/internal/repository"
	"clash-server/pkg/crypto"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService() *AuthService {
	return &AuthService{userRepo: repository.NewUserRepository()}
}

func (s *AuthService) Login(username, password string) (*model.User, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if !crypto.CheckPassword(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

func (s *AuthService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if !crypto.CheckPassword(oldPassword, user.Password) {
		return errors.New("invalid old password")
	}
	hashedPassword, err := crypto.HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return s.userRepo.Update(user)
}

func (s *AuthService) InitPassword(username, password string) error {
	count, err := s.userRepo.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("user already exists")
	}
	hashedPassword, err := crypto.HashPassword(password)
	if err != nil {
		return err
	}
	user := &model.User{
		Username: username,
		Password: hashedPassword,
	}
	return s.userRepo.Create(user)
}

func (s *AuthService) IsInitialized() (bool, error) {
	count, err := s.userRepo.Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
