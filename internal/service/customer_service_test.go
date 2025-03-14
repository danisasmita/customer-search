package service_test

import (
	"errors"
	"testing"

	"github.com/danisasmita/customer-search/internal/model"
	"github.com/danisasmita/customer-search/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	customerName1  = "A John Doe"
	customerEmail1 = "ajohn.doe@example.com"
	accountNumber1 = "123456"

	customerName2  = "Jane Doe"
	customerEmail2 = "jane.doe@example.com"
	accountNumber2 = "654321"
)

// MockCustomerRepository adalah mock untuk CustomerRepository
type MockCustomerRepository struct {
	mock.Mock
}

// FindByName adalah metode mock untuk mencari customer berdasarkan nama, email, dan nomor akun
func (m *MockCustomerRepository) FindByName(name, email, accountNumber string) ([]model.Customer, error) {
	args := m.Called(name, email, accountNumber)
	return args.Get(0).([]model.Customer), args.Error(1)
}

func TestCustomerServiceSearchByName(t *testing.T) {
	// Buat instance mock repository
	mockRepo := new(MockCustomerRepository)
	// Buat instance CustomerService dengan mock repository
	customerService := service.NewCustomerService(mockRepo)

	// Data dummy untuk customer
	dummyCustomers := []model.Customer{
		{
			Name:  customerName1,
			Email: customerEmail1,
			BankAccounts: []model.BankAccount{
				{AccountNumber: accountNumber1, Balance: 1000.0},
			},
		},
		{
			Name:  customerName2,
			Email: customerEmail2,
			BankAccounts: []model.BankAccount{
				{AccountNumber: accountNumber2, Balance: 2000.0},
			},
		},
	}

	t.Run("success - found customers", func(t *testing.T) {
		// Set up mock behavior
		mockRepo.On("FindByName", customerName1, customerEmail1, accountNumber1).
			Return(dummyCustomers, nil).
			Once()

		// Panggil method yang di-test
		result, err := customerService.SearchByName(customerName1, customerEmail1, accountNumber1)

		// Assert hasil
		assert.NoError(t, err)
		assert.Equal(t, dummyCustomers, result)
		// Pastikan mock dipanggil sesuai ekspektasi
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - repository error", func(t *testing.T) {
		// Set up mock behavior untuk mengembalikan error
		mockRepo.On("FindByName", "John Doe", "john.doe@example.com", "123456").
			Return([]model.Customer{}, errors.New("database error")).
			Once()

		// Panggil method yang di-test
		result, err := customerService.SearchByName("John Doe", "john.doe@example.com", "123456")

		// Assert hasil
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		assert.Empty(t, result)
		// Pastikan mock dipanggil sesuai ekspektasi
		mockRepo.AssertExpectations(t)
	})

	t.Run("success - no customers found", func(t *testing.T) {
		// Set up mock behavior untuk mengembalikan slice kosong
		mockRepo.On("FindByName", "Unknown", "", "").
			Return([]model.Customer{}, nil).
			Once()

		// Panggil method yang di-test
		result, err := customerService.SearchByName("Unknown", "", "")

		// Assert hasil
		assert.NoError(t, err)
		assert.Empty(t, result)
		// Pastikan mock dipanggil sesuai ekspektasi
		mockRepo.AssertExpectations(t)
	})
}
