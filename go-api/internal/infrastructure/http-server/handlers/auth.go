package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server/contracts"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server/converter"
	"github.com/pttrulez/investor-go-next/go-api/internal/service/user"
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

	// Валидация пришедших данных
	if err := c.validator.Struct(loginData); err != nil {
		var validateErr validator.ValidationErrors
		errors.As(err, &validateErr)
		writeValidationErrorsJSON(w, validateErr)
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
		"id":         u.ID,
		"email":      u.Email,
		"name":       u.Name,
		"role":       u.Role,
		"tg_chat_id": u.TgChatID,
	}

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

	// Валидация пришедших данных
	if err := c.validator.Struct(registerData); err != nil {
		var validateErr validator.ValidationErrors
		errors.As(err, &validateErr)
		writeValidationErrorsJSON(w, validateErr)
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
