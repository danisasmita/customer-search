package service

import (
	"github.com/danisasmita/customer-search/internal/model"
	"github.com/danisasmita/customer-search/internal/repository"
)

type CustomerService struct {
	repo repository.CustomerRepository
}

func NewCustomerService(repo repository.CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

func (s *CustomerService) SearchByName(name, email, accountNumber string) ([]model.Customer, error) {
	return s.repo.FindByName(name, email, accountNumber)
}
