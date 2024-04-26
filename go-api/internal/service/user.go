package service

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pttrulez/investor-go/internal/model"
	"github.com/pttrulez/investor-go/internal/repository"
	httpresponse "github.com/pttrulez/investor-go/pkg/http-response"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	tokenAuth *jwtauth.JWTAuth
	repo      *repository.Repository
}

func (s *UserService) LoginUser(ctx context.Context, user *model.User) (string, error) {
	user, err := s.repo.User.GetByEmail(ctx, user.Email)
	if err != nil {
		return "", err
	} else if user == nil {
		return "", httpresponse.NewErrSendToClient("Пользователя с таким email не существует", http.StatusBadRequest)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(user.Password))
	if err != nil {
		return "", httpresponse.NewErrSendToClient("Неверные данные", http.StatusUnauthorized)
	}

	claims := jwt.MapClaims{
		"id":    user.Id,
		"email": user.Email,
		"name":  user.Name,
		"role":  user.Role,
	}
	jwtauth.SetExpiry(claims, time.Now().Add(time.Hour*6))
	_, tokenString, _ := s.tokenAuth.Encode(claims)

	return tokenString, nil
}

func (s *UserService) RegisterUser(ctx context.Context, user *model.User) error {
	// Check if user with this email already exists
	existingUser, err := s.repo.User.GetByEmail(ctx, user.Email)
	if existingUser != nil {
		return httpresponse.NewErrSendToClient("Пользователь с таким email уже существует", http.StatusBadRequest)
	} else if err != nil {
		return err
	}

	encpw, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.HashedPassword = string(encpw)
	// Creating new user
	err = s.repo.User.Insert(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func NewUserService(repo *repository.Repository) *UserService {
	return &UserService{
		repo: repo,
	}
}
