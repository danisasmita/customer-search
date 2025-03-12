package service

import (
	"errors"

	"github.com/danisasmita/customer-search/internal/model"
	"github.com/danisasmita/customer-search/internal/repository"
	"github.com/danisasmita/customer-search/pkg/message"
	"github.com/danisasmita/customer-search/pkg/utils"
)

type AuthService interface {
	Register(user *model.User) error
	Login(request model.UserRequest) (string, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) Register(user *model.User) error {
	if user.Username == "" || user.Password == "" {
		return errors.New(message.UsernameRequired + " and " + message.PasswordRequired)
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return errors.New(message.InternalServerError)
	}
	user.Password = hashedPassword
	return s.repo.CreateUser(user)
}

func (s *authService) Login(request model.UserRequest) (string, error) {
	if request.Username == "" || request.Password == "" {
		return "", errors.New(message.UsernameRequired + " and " + message.PasswordRequired)
	}

	user, err := s.repo.FindUserByUsername(request.Username)
	if err != nil {
		return "", errors.New(message.UserNotFound)
	}

	if !utils.CheckPasswordHash(request.Password, user.Password) {
		return "", errors.New(message.InvalidCredentials)
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return "", errors.New(message.InternalServerError)
	}

	return token, nil
}
