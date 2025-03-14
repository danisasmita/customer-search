package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danisasmita/customer-search/internal/handler"
	"github.com/danisasmita/customer-search/internal/model"
	"github.com/danisasmita/customer-search/internal/service"
	"github.com/danisasmita/customer-search/pkg/message"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock service untuk AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockAuthService) Login(user model.UserRequest) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

// Tambahkan setelah MockAuthService
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindUserByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	return args.Get(0).(*model.User), args.Error(1)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.Default()
}

const (
	registerPath      = "/register"
	loginPath         = "/login"
	contentType       = "application/json"
	contentTypeHeader = "Content-Type"
)

func TestAuthHandlerRegister(t *testing.T) {
	mockService := new(MockAuthService)
	authHandler := handler.NewAuthHandler(mockService)
	router := setupRouter()
	router.POST(registerPath, authHandler.Register)

	t.Run("success - register", func(t *testing.T) {
		user := model.User{Username: "john_doe", Password: "password123"}
		mockService.On("Register", &user).Return(nil)

		reqBody, _ := json.Marshal(user)
		req, _ := http.NewRequest(http.MethodPost, registerPath, bytes.NewBuffer(reqBody))
		req.Header.Set(contentTypeHeader, contentType)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("error - missing username or password", func(t *testing.T) {
		user := model.User{}
		reqBody, _ := json.Marshal(user)
		req, _ := http.NewRequest(http.MethodPost, registerPath, bytes.NewBuffer(reqBody))
		req.Header.Set(contentTypeHeader, contentType)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("error - internal server error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockService := service.NewAuthService(mockRepo)
		authHandler := handler.NewAuthHandler(mockService)

		router := setupRouter()
		router.POST(registerPath, authHandler.Register)

		user := model.User{Username: "john_doe", Password: "password123"}

		// Mocking `CreateUser` untuk mengembalikan error
		mockRepo.On("CreateUser", mock.AnythingOfType("*model.User")).Return(errors.New("internal error")).Once()

		reqBody, _ := json.Marshal(user)
		req, _ := http.NewRequest(http.MethodPost, registerPath, bytes.NewBuffer(reqBody))
		req.Header.Set(contentTypeHeader, contentType)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		mockRepo.AssertExpectations(t)
	})

}

func TestAuthHandlerLogin(t *testing.T) {
	mockService := new(MockAuthService)
	authHandler := handler.NewAuthHandler(mockService)
	router := setupRouter()
	router.POST(loginPath, authHandler.Login)

	t.Run("success - login", func(t *testing.T) {
		user := model.UserRequest{Username: "john_doe", Password: "password123"}
		mockService.On("Login", user).Return("valid-token", nil)

		reqBody, _ := json.Marshal(user)
		req, _ := http.NewRequest(http.MethodPost, loginPath, bytes.NewBuffer(reqBody))
		req.Header.Set(contentTypeHeader, contentType)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("error - missing username or password", func(t *testing.T) {
		user := model.UserRequest{}
		reqBody, _ := json.Marshal(user)
		req, _ := http.NewRequest(http.MethodPost, loginPath, bytes.NewBuffer(reqBody))
		req.Header.Set(contentTypeHeader, contentType)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("error - invalid credentials", func(t *testing.T) {
		user := model.UserRequest{Username: "john_doe", Password: "wrongpassword"}
		mockService.On("Login", user).Return("", errors.New(message.InvalidCredentials))

		reqBody, _ := json.Marshal(user)
		req, _ := http.NewRequest(http.MethodPost, loginPath, bytes.NewBuffer(reqBody))
		req.Header.Set(contentTypeHeader, contentType)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		mockService.AssertExpectations(t)
	})
}
