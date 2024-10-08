package converter

import (
	"github.com/pttrulez/investor-go/internal/domain"
	"github.com/pttrulez/investor-go/internal/infrastructure/http-server/contracts"
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

func FromUpdateUserRequestToUser(dto contracts.UpdateUserRequest) domain.User {
	var res domain.User

	if dto.Name != nil {
		res.Name = *dto.Name
	}
	res.TgChatID = dto.InvestBotTgChatId

	return res
}
