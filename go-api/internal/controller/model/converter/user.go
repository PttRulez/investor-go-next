package converter

import (
	"github.com/pttrulez/investor-go/internal/controller/model/dto"
	"github.com/pttrulez/investor-go/internal/entity"
)

func FromRegisterDataToUser(dto *dto.RegisterUser) *entity.User {
	user := entity.User{
		Email:    dto.Email,
		Name:     dto.Name,
		Password: dto.Password,
	}

	if dto.Role == "" {
		user.Role = user.Investor
	} else {
		user.Role = dto.Role
	}

	return &user
}

func FromLoginDtoToUser(dto *dto.LoginUser) *entity.User {
	return &entity.User{
		Email:    dto.Email,
		Password: dto.Password,
	}
}
