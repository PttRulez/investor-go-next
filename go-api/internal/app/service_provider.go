package app

import (
	"github.com/pttrulez/investor-go/internal/repository"
	"github.com/pttrulez/investor-go/internal/service"
)

func NewServiceContainer(repo *repository.Repository) *service.Container {
	container := service.Container{}
	container.User = service.NewUserService(repo)
	container.Validator = service.NewValidator()

	return &container
}
