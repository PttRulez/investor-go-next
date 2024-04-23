package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pttrulez/investor-go/internal/lib/response"
	"github.com/pttrulez/investor-go/internal/types"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	repo      *types.Repository
	tokenAuth *jwtauth.JWTAuth
}

func (c *AuthController) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var loginData types.LoginUser

	// Return error if json couldn't be decoded
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := c.repo.User.GetByEmail(ctx, loginData.Email)
	fmt.Println(loginData.Email)
	if err != nil {
		response.WriteString(w, http.StatusBadRequest, err.Error())
		return
	} else if user == nil {
		response.WriteString(w, http.StatusBadRequest, "Пользователя с таким email не существует")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(loginData.Password))
	if err != nil {
		response.WriteString(w, http.StatusUnauthorized, "Неверные данные")
		return
	}

	claims := jwt.MapClaims{
		"id":    user.Id,
		"email": user.Email,
		"name":  user.Name,
		"role":  user.Role,
	}
	jwtauth.SetExpiry(claims, time.Now().Add(time.Hour*6))
	_, tokenString, _ := c.tokenAuth.Encode(claims)

	response.WriteOKJSON(w, map[string]string{"token": tokenString})
}

func (c *AuthController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var registerData types.RegisterUser

	// Return error if json couldn't be decoded
	if err := json.NewDecoder(r.Body).Decode(&registerData); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate request fields
	if err := validator.New().Struct(registerData); err != nil {
		validateErr := err.(validator.ValidationErrors)
		response.WriteValidationErrorsJSON(w, validateErr)
		return
	}

	// Check if user with this email already exists
	existingUser, err := c.repo.User.GetByEmail(ctx, registerData.Email)
	if existingUser != nil {
		response.WriteString(w, http.StatusBadRequest, "Пользователь с таким email уже существует")
		return
	} else if err != nil {
		response.WriteString(w, http.StatusBadRequest, err.Error())
		return
	}

	encpw, err := bcrypt.GenerateFromPassword([]byte(registerData.Password), bcrypt.DefaultCost)
	if err != nil {
		response.WriteString(w, http.StatusInternalServerError, ":(")
	}

	if registerData.Role == "" {
		registerData.Role = types.Investor
	}

	// Creating new user
	err = c.repo.User.Insert(ctx, types.User{
		Email:          registerData.Email,
		Name:           registerData.Name,
		HashedPassword: string(encpw),
		Role:           registerData.Role,
	})
	if err != nil {
		response.WriteString(w, http.StatusInternalServerError, "Failed to create new user")
	}

	w.WriteHeader(http.StatusCreated)
}

func NewAuthController(repo *types.Repository, tokenAuth *jwtauth.JWTAuth) *AuthController {
	return &AuthController{
		repo:      repo,
		tokenAuth: tokenAuth,
	}
}
