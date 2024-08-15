package httpcontrollers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/pttrulez/investor-go/internal/controller/model/converter"
	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/service/user"
	"github.com/pttrulez/investor-go/pkg/api"
)

func (c *AuthController) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var loginData api.LoginRequest

	// Return error if json couldn't be decoded
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	tokenString, err := c.userService.LoginUser(ctx, converter.FromLoginRequestToUser(loginData))
	if errors.Is(err, user.ErrWrongUsername) || errors.Is(err, user.ErrWrongPassword) {
		writeString(w, http.StatusUnauthorized, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}

func (c *AuthController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var registerData api.RegisterUserRequest

	// Return error if json couldn't be decoded
	if err := json.NewDecoder(r.Body).Decode(&registerData); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate request fields
	if err := validator.New().Struct(registerData); err != nil {
		var validateErr validator.ValidationErrors
		if ok := errors.As(err, &validateErr); ok {
			writeValidationErrorsJSON(w, validateErr)
		} else {
			writeError(w, err)
		}
		return
	}

	err := c.userService.RegisterUser(ctx, converter.FromRegisterRequestToUser(registerData))
	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type UserService interface {
	LoginUser(ctx context.Context, user *entity.User) (string, error)
	RegisterUser(ctx context.Context, user *entity.User) error
}

type AuthController struct {
	userService UserService
	validator   *validator.Validate
}

func NewAuthController(userService UserService, validator *validator.Validate) *AuthController {
	return &AuthController{
		userService: userService,
		validator:   validator,
	}
}
