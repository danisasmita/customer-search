package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danisasmita/customer-search/internal/handler"
	"github.com/danisasmita/customer-search/internal/model"
	"github.com/danisasmita/customer-search/pkg/message"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	searchEndpoint = "/search?name=John"
)

// MockCustomerService adalah mock untuk service.CustomerService
type MockCustomerService struct {
	mock.Mock
}

func (m *MockCustomerService) SearchByName(name, email, accountNumber string) ([]model.Customer, error) {
	args := m.Called(name, email, accountNumber)
	return args.Get(0).([]model.Customer), args.Error(1)
}

func setRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.Default()
}

// Tambahkan setelah MockAuthService
type MockCustomerRepository struct {
	mock.Mock
}

func TestCustomerHandlerSearchByName(t *testing.T) {
	mockService := new(MockCustomerService)
	customerHandler := handler.NewCustomerHandler(mockService)
	router := setupRouter()
	router.GET("/search", customerHandler.SearchByName)

	t.Run("success - search by name", func(t *testing.T) {
		// Mock data
		mockCustomers := []model.Customer{
			{
				Name:  "John Doe",
				Email: "john@example.com",
				BankAccounts: []model.BankAccount{
					{AccountNumber: "123456"},
				},
			},
		}
		mockService.On("SearchByName", "John", "", "").Return(mockCustomers, nil)

		// Buat request
		req, _ := http.NewRequest(http.MethodGet, searchEndpoint, nil)
		recorder := httptest.NewRecorder()

		// Jalankan request
		router.ServeHTTP(recorder, req)

		// Assertions
		assert.Equal(t, http.StatusOK, recorder.Code)

		// Unmarshal response ke struct yang sesuai
		var response struct {
			Data []model.Customer `json:"data"`
		}
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Bandingkan data
		assert.Equal(t, mockCustomers, response.Data)

		mockService.AssertExpectations(t)
	})

	t.Run("error - no query parameters provided", func(t *testing.T) {
		// Buat request tanpa query parameters
		req, _ := http.NewRequest(http.MethodGet, "/search", nil)
		recorder := httptest.NewRecorder()

		// Jalankan request
		router.ServeHTTP(recorder, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		var responseBody gin.H
		err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
		assert.NoError(t, err)
		assert.Equal(t, gin.H{"error": message.SearchCustomer}, responseBody)
	})
}
