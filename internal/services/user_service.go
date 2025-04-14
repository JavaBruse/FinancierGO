package services

import (
	"errors"
	"financierGo/internal/models"
	"financierGo/internal/repositories"
	"financierGo/internal/utils"
)

type UserService struct {
	Repo *repositories.UserRepository
}

func (s *UserService) Register(username, email, password string) (*models.User, error) {
	existingUser, _ := s.Repo.FindByEmail(email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	hashedPwd, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Password: hashedPwd,
	}

	err = s.Repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(email, password string) (*models.User, error) {
	user, err := s.Repo.FindByEmail(email)
	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
