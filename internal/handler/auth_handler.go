package handler

import (
	"fmt"
	"net/http"

	"errors"

	"github.com/danisasmita/customer-search/internal/model"
	"github.com/danisasmita/customer-search/internal/service"
	"github.com/danisasmita/customer-search/pkg/message"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": message.BadRequest})
		return
	}

	if user.Username == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": message.UsernameRequired + " and " + message.PasswordRequired})
		return
	}

	if err := h.service.Register(&user); err != nil {
		fmt.Println("Error registering user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": message.InternalServerError})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": message.UserRegistered})
}

func (h *AuthHandler) ValidateCredentials(user model.UserRequest) error {
	if user.Username == "" || user.Password == "" {
		return errors.New(message.UsernameRequired + " and " + message.PasswordRequired)
	}
	return nil
}

func (h *AuthHandler) Login(c *gin.Context) {
	var user model.UserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": message.BadRequest})
		return
	}

	if err := h.ValidateCredentials(user); err != nil {
		fmt.Printf("Validation error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(user)
	if err != nil {
		fmt.Printf("Login error: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": message.InvalidCredentials})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
