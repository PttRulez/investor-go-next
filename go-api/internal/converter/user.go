package converter

import (
	"github.com/pttrulez/investor-go/internal/api/model/dto"
	"github.com/pttrulez/investor-go/internal/model"
)

func FromRegisterDataToUser(dto *dto.RegisterUser) *model.User {
	user := model.User{
		Email:    dto.Email,
		Name:     dto.Name,
		Password: dto.Password,
	}

	if dto.Role == "" {
		user.Role = model.Investor
	} else {
		user.Role = dto.Role
	}

	return &user
}

func FromLoginDtoToUser(dto *dto.LoginUser) *model.User {
	return &model.User{
		Email:    dto.Email,
		Password: dto.Password,
	}
}
