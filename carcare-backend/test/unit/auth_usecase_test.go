package unit

import (
	"context"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/VladimirGrebenev/CarCare-backend/internal/usecase"
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/user"
)

type mockUserRepo struct {
	users map[string]*user.User
}

func (m *mockUserRepo) Create(ctx context.Context, u *user.User) error {
	if _, exists := m.users[u.Email]; exists {
		return user.ErrAlreadyExists
	}
	m.users[u.Email] = u
	return nil
}
func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	u, ok := m.users[email]
	if !ok {
		return nil, user.ErrNotFound
	}
	return u, nil
}
// ...другие методы user.Repository

func TestRegister_Success(t *testing.T) {
	repo := &mockUserRepo{users: map[string]*user.User{}}
	uc := usecase.NewAuthUsecase(repo)
	err := uc.Register(context.Background(), "test@example.com", "password")
	assert.NoError(t, err)
}

func TestRegister_Duplicate(t *testing.T) {
	repo := &mockUserRepo{users: map[string]*user.User{"test@example.com": {Email: "test@example.com"}}}
	uc := usecase.NewAuthUsecase(repo)
	err := uc.Register(context.Background(), "test@example.com", "password")
	assert.Error(t, err)
}

// ...тесты ConfirmEmail, Login, ForgotPassword, ResetPassword, OAuth, RefreshToken, Logout
