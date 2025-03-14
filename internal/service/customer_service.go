package service

import (
	"github.com/danisasmita/customer-search/internal/model"
	"github.com/danisasmita/customer-search/internal/repository"
)

type CustomerService interface {
	SearchByName(name, email, accountNumber string) ([]model.Customer, error)
}

// CustomerServiceImpl adalah implementasi dari CustomerService
type CustomerServiceImpl struct {
	repo repository.CustomerRepository
}

// NewCustomerService menginisialisasi CustomerServiceImpl
func NewCustomerService(repo repository.CustomerRepository) *CustomerServiceImpl {
	return &CustomerServiceImpl{repo: repo}
}

// SearchByName mencari pelanggan berdasarkan nama, email, dan nomor akun
func (s *CustomerServiceImpl) SearchByName(name, email, accountNumber string) ([]model.Customer, error) {
	customers, err := s.repo.FindByName(name, email, accountNumber)

	if err != nil {
		return nil, err
	}

	// Jika data tidak ditemukan, kembalikan slice kosong
	if len(customers) == 0 {
		return []model.Customer{}, nil
	}

	return customers, nil
}
