package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/pttrulez/investor-go/internal/domain"
	"github.com/pttrulez/investor-go/internal/infrastructure/database"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrWrongUsername = errors.New("неправильный юзернейм")
	ErrWrongPassword = errors.New("неправильный пароль")
)



func (s *Service) LoginUser(ctx context.Context, model domain.User) (domain.User, error) {
	const op = "UserService.LoginUser"

	user, err := s.userRepo.GetByEmail(ctx, model.Email)
	if errors.Is(err, database.ErrNotFound) {
		return domain.User{}, ErrWrongUsername
	}
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(model.Password))
	if err != nil {
		return domain.User{}, ErrWrongPassword
	}

	return user, nil
}

func (s *Service) RegisterUser(ctx context.Context, user domain.User) error {
	const op = "UserService.RegisterUser"

	// Check if user with this email already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, user.Email)
	if existingUser.ID != 0 {
		return errors.New("такой юзер уже существует")
	}
	if err != nil && !errors.Is(err, database.ErrNotFound) {
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
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

type Repository interface {
	Insert(ctx context.Context, u domain.User) error
	GetByEmail(ctx context.Context, email string) (domain.User, error)
}

type Service struct {
	userRepo Repository
}

func NewUserService(repo Repository) *Service {
	return &Service{
		userRepo: repo,
	}
}
