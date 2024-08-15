package converter

import (
	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/pkg/api"
)

func FromRegisterRequestToUser(dto api.RegisterUserRequest) *entity.User {
	user := entity.User{
		Email:    dto.Email,
		Name:     dto.Name,
		Password: dto.Password,
	}

	if dto.Role == nil {
		user.Role = entity.Investor
	} else {
		var role entity.Role
		switch *dto.Role {
		case api.ADMIN:
			role = entity.Admin
		case api.INVESTOR:
			role = entity.Investor
		default:
			role = entity.Investor
		}
		user.Role = role
	}

	return &user
}

func FromLoginRequestToUser(dto api.LoginRequest) *entity.User {
	return &entity.User{
		Email:    dto.Email,
		Password: dto.Password,
	}
}
