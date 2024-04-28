package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/pttrulez/investor-go/internal/api/model/converter"
	"github.com/pttrulez/investor-go/internal/api/model/dto"
	"github.com/pttrulez/investor-go/internal/model"
	"github.com/pttrulez/investor-go/internal/repository"
	"github.com/pttrulez/investor-go/internal/service"
	httpresponse "github.com/pttrulez/investor-go/pkg/http-response"
)

func (c *AuthController) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var loginData dto.LoginUser

	// Return error if json couldn't be decoded
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		httpresponse.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Printf("loginData: %v\n", loginData)

	tokenString, err := c.services.User.LoginUser(ctx, converter.FromLoginDtoToUser(&loginData))
	if err != nil {
		httpresponse.SendError(w, err)
		return
	}
	fmt.Println("IM HERE")
	httpresponse.WriteOKJSON(w, map[string]string{"token": tokenString})
}

func (c *AuthController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var registerData dto.RegisterUser

	// Return error if json couldn't be decoded
	if err := json.NewDecoder(r.Body).Decode(&registerData); err != nil {
		httpresponse.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate request fields
	if err := validator.New().Struct(registerData); err != nil {
		validateErr := err.(validator.ValidationErrors)
		httpresponse.WriteValidationErrorsJSON(w, validateErr)
		return
	}

	err := c.services.User.RegisterUser(ctx, converter.FromRegisterDataToUser(&registerData))
	if err != nil {
		httpresponse.SendError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type UserService interface {
	LoginUser(ctx context.Context, user *model.User) (string, error)
	RegisterUser(ctx context.Context, user *model.User) error
}

type AuthController struct {
	repo     *repository.Repository
	services *service.Container
}

func NewAuthController(repo *repository.Repository, services *service.Container) *AuthController {
	return &AuthController{
		repo:     repo,
		services: services,
	}
}
