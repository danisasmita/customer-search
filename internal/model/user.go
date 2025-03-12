package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

type (
	UserRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
)

type (
	LoginResponse struct {
		AccessToken string `json:"accessToken"`
	}
)
