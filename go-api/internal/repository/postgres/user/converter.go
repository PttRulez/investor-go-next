package user

import (
	"github.com/pttrulez/investor-go/internal/model"
)

func FromUserToDBUser(user *model.User) *User {
	return &User{
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
		Id:             user.Id,
		Name:           user.Name,
		Role:           string(user.Role),
	}
}

func FromDBToModelUser(user *User) *model.User {
	return &model.User{
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
		Id:             user.Id,
		Name:           user.Name,
		Role:           model.Role(user.Role),
	}
}
