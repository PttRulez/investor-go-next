package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pttrulez/investor-go/internal/api/contracts"
	"github.com/pttrulez/investor-go/internal/api/converter"
	"github.com/pttrulez/investor-go/internal/service/user"
)

const tokentExpHours = 6

func (c *Handlers) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var loginData contracts.LoginRequest

	// Return error if json couldn't be decoded
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	u, err := c.userService.LoginUser(ctx, converter.FromLoginRequestToUser(loginData))
	if errors.Is(err, user.ErrWrongUsername) || errors.Is(err, user.ErrWrongPassword) {
		writeString(w, http.StatusUnauthorized, err.Error())
		return
	}
	if err != nil {
		c.logger.Error(err)
		writeString(w, http.StatusInternalServerError, "Не удалось авторизоваться")
		return
	}

	claims := jwt.MapClaims{
		"id":    u.ID,
		"email": u.Email,
		"name":  u.Name,
		"role":  u.Role,
	}

	fmt.Println("CLAIMS:", claims)

	jwtauth.SetExpiry(claims, time.Now().Add(time.Hour*tokentExpHours))

	_, tokenString, err := c.tokenAuth.Encode(claims)
	if err != nil {
		c.logger.Error(err)
		writeString(w, http.StatusInternalServerError, "Не удалось авторизоваться")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}

func (c *Handlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var registerData contracts.RegisterUserRequest

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
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// type AuthController struct {
// 	logger      *logger.Logger
// 	tokenAuth   *jwtauth.JWTAuth
// 	userService interfaces.UserService
// 	validator   *validator.Validate
// }

// func NewAuthController(logger *logger.Logger, tokenAuth *jwtauth.JWTAuth,
// 	userService interfaces.UserService, validator *validator.Validate) *AuthController {
// 	return &AuthController{
// 		logger:      logger,
// 		tokenAuth:   tokenAuth,
// 		userService: userService,
// 		validator:   validator,
// 	}
// }
