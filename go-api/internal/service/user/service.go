package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/infrastracture/database"

	"github.com/go-chi/jwtauth/v5"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	errWrongUsername = errors.New("неправильный юзернейм")
	errWrongPassword = errors.New("неправильный юзернейм")
)

const tokentExpHours = 6

func (s *Service) LoginUser(ctx context.Context, model *entity.User) (string, error) {
	const op = "UserService.LoginUser"

	user, err := s.userRepo.GetByEmail(ctx, model.Email)
	if errors.Is(err, database.ErrNotFound) {
		return "", errWrongUsername
	}
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(model.Password))
	if err != nil {
		return "", errWrongPassword
	}

	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
		"role":  user.Role,
	}

	jwtauth.SetExpiry(claims, time.Now().Add(time.Hour*tokentExpHours))

	_, tokenString, err := s.tokenAuth.Encode(claims)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Service) RegisterUser(ctx context.Context, user *entity.User) error {
	const op = "UserService.LoginUser"

	// Check if user with this email already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, user.Email)
	if existingUser != nil {
		return errors.New("такой юзер уже существует")
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
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
