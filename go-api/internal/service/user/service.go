package user

import (
	"context"
	"github.com/pttrulez/investor-go/internal/entity"
	ierrors "github.com/pttrulez/investor-go/internal/errors"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) LoginUser(ctx context.Context, model *entity.User) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, model.Email)
	if err != nil {
		return "", err
	} else if user == nil {
		return "", ierrors.NewErrSendToClient("Неверные данные", http.StatusUnauthorized)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(model.Password))
	if err != nil {
		return "", ierrors.NewErrSendToClient("Неверные данные", http.StatusUnauthorized)
	}

	claims := jwt.MapClaims{
		"id":    user.Id,
		"email": user.Email,
		"name":  user.Name,
		"role":  user.Role,
	}
	jwtauth.SetExpiry(claims, time.Now().Add(time.Hour*6))

	_, tokenString, err := s.tokenAuth.Encode(claims)
	if err != nil {
		return "", nil
	}

	return tokenString, nil
}

func (s *Service) RegisterUser(ctx context.Context, user *entity.User) error {
	// Check if user with this email already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, user.Email)
	if existingUser != nil {
		return ierrors.NewErrSendToClient("Пользователь с таким email уже существует", http.StatusBadRequest)
	} else if err != nil {
		return err
	}

	encpw, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.HashedPassword = string(encpw)
	// Creating new user
	err = s.userRepo.Insert(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

type Repository interface {
	Insert(ctx context.Context, u *entity.User) error
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}

type Service struct {
	userRepo  Repository
	tokenAuth *jwtauth.JWTAuth
}

func NewUserService(repo Repository, tokenAuth *jwtauth.JWTAuth) *Service {
	return &Service{
		tokenAuth: tokenAuth,
		userRepo:  repo,
	}
}
