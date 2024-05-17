package dto

import (
	"github.com/pttrulez/investor-go/internal/entity"
)

type RegisterUser struct {
	Email    string      `json:"email" validate:"required,email"`
	Name     string      `json:"name" validate:"required"`
	Password string      `json:"password" validate:"required"`
	Role     entity.Role `json:"role"`
}

type LoginUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
