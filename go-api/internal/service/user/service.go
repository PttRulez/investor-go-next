package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/storage"
	"github.com/pttrulez/investor-go-next/go-api/internal/service"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrWrongUsername = errors.New("неправильный юзернейм")
	ErrWrongPassword = errors.New("неправильный пароль")
)

func (s *Service) GetUserByChatID(ctx context.Context, chatID string) (domain.User, error) {
	const op = "UserService.GetUserByChatID"

	user, err := s.repo.GetUserByChatID(ctx, chatID)
	if errors.Is(err, storage.ErrNotFound) {
		return domain.User{}, service.ErrDomainNotFound
	}
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Service) LoginUser(ctx context.Context, model domain.User) (domain.User, error) {
	const op = "UserService.LoginUser"

	user, err := s.repo.GetUserByEmail(ctx, model.Email)
	if errors.Is(err, storage.ErrNotFound) {
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
	existingUser, err := s.repo.GetUserByEmail(ctx, user.Email)
	if existingUser.ID != 0 {
		return errors.New("такой юзер уже существует")
	}
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return fmt.Errorf("%s: %w", op, err)
	}

	encpw, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.HashedPassword = string(encpw)
	// Creating new user
	err = s.repo.InsertUser(ctx, user)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) UpdateUser(ctx context.Context, user domain.User) error {
	const op = "UserService.UpdateUser"

	err := s.repo.UpdateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	GetUserByChatID(ctx context.Context, chatID string) (domain.User, error)
	InsertUser(ctx context.Context, u domain.User) error
	UpdateUser(ctx context.Context, u domain.User) error
}

type Service struct {
	repo Repository
}

func NewUserService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}
