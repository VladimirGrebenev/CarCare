package usecase

import (
	"context"
	"errors"
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/user"
	"github.com/google/uuid"
)

type UserService interface {
	Create(ctx context.Context, u *user.User) error
	Get(ctx context.Context, id string) (*user.User, error)
	Update(ctx context.Context, u *user.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*user.User, error)
}

type userService struct {
	repo user.Repository
}

func NewUserService(repo user.Repository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(ctx context.Context, u *user.User) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	if err := u.Validate(); err != nil {
		return err
	}
	// Проверка уникальности email
	existing, _ := s.repo.GetByEmail(ctx, u.Email)
	if existing != nil {
		return errors.New("user with this email already exists")
	}
	return s.repo.Create(ctx, u)
}

func (s *userService) Get(ctx context.Context, id string) (*user.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) Update(ctx context.Context, u *user.User) error {
	if err := u.Validate(); err != nil {
		return err
	}
	return s.repo.Update(ctx, u)
}

func (s *userService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *userService) List(ctx context.Context) ([]*user.User, error) {
	return s.repo.List(ctx)
}
