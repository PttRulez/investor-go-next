package converter

import (
	"github.com/pttrulez/investor-go/internal/api/contracts"
	"github.com/pttrulez/investor-go/internal/domain"
)

func FromRegisterRequestToUser(dto contracts.RegisterUserRequest) domain.User {
	user := domain.User{
		Email:    dto.Email,
		Name:     dto.Name,
		Password: dto.Password,
	}

	if dto.Role == nil {
		user.Role = domain.Investor
	} else {
		var role domain.Role
		switch *dto.Role {
		case contracts.ADMIN:
			role = domain.Admin
		case contracts.INVESTOR:
			role = domain.Investor
		default:
			role = domain.Investor
		}
		user.Role = role
	}

	return user
}

func FromLoginRequestToUser(dto contracts.LoginRequest) domain.User {
	return domain.User{
		Email:    dto.Email,
		Password: dto.Password,
	}
}
