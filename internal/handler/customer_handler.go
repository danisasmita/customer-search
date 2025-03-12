package handler

import (
	"net/http"

	"github.com/danisasmita/customer-search/internal/service"
	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	service service.CustomerService
}

func NewCustomerHandler(service service.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

func (h *CustomerHandler) SearchByName(c *gin.Context) {
	name := c.Query("name")
	email := c.Query("email")
	accountNumber := c.Query("account_number")

	if name == "" && email == "" && accountNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Harap berikan setidaknya name, email, atau account_number untuk pencarian"})
		return
	}

	customers, err := h.service.SearchByName(name, email, accountNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan pada server"})
		return
	}

	if len(customers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Customers found",
		"data":    customers,
	})
}
