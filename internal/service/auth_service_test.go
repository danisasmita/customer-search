package service_test

import (
	"errors"
	"testing"

	"github.com/danisasmita/customer-search/internal/model"
	"github.com/danisasmita/customer-search/internal/service"
	"github.com/danisasmita/customer-search/pkg/message"
	"github.com/danisasmita/customer-search/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockUserRepository adalah mock untuk UserRepository
type MockUserRepository struct {
	mock.Mock
}

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// Mock metode FindUserByUsername
func (m *MockAuthService) Login(request model.UserRequest) (string, error) {
	args := m.Called(request)
	return args.String(0), args.Error(1)
}

func TestAuthServiceRegister(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := service.NewAuthService(mockRepo)

	t.Run("success - register user", func(t *testing.T) {
		user := &model.User{Username: "john_doe", Password: "password123"}
		mockRepo.On("CreateUser", user).Return(nil)

		err := authService.Register(user)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - missing username or password", func(t *testing.T) {
		user := &model.User{}
		err := authService.Register(user)
		assert.EqualError(t, err, message.UsernameRequired+" and "+message.PasswordRequired)
	})
}

// FindUserByUsername adalah metode mock untuk menemukan pengguna berdasarkan username
func (m *MockUserRepository) FindUserByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) CreateUser(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestAuthServiceLogin(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := service.NewAuthService(mockRepo)

	// Hash password untuk simulasi user di database
	hashedPassword, err := utils.HashPassword("password123")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	validUser := &model.User{
		Model:    gorm.Model{ID: 1},
		Username: "john_doe",
		Password: hashedPassword,
	}

	// Mocking repository behavior
	mockRepo.On("FindUserByUsername", "john_doe").Return(validUser, nil)
	mockRepo.On("FindUserByUsername", "unknown").Return((*model.User)(nil), errors.New("not found"))

	t.Run("success - login", func(t *testing.T) {
		request := model.UserRequest{Username: "john_doe", Password: "password123"}
		token, err := authService.Login(request)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("error - user not found", func(t *testing.T) {
		request := model.UserRequest{Username: "unknown", Password: "password123"}
		_, err := authService.Login(request)
		assert.EqualError(t, err, message.UserNotFound)
	})

	t.Run("error - wrong password", func(t *testing.T) {
		request := model.UserRequest{Username: "john_doe", Password: "wrongpassword"}
		_, err := authService.Login(request)
		assert.EqualError(t, err, message.InvalidCredentials)
	})
}
