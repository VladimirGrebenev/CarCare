package unit

import (
	"context"
	"testing"

	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/user"
	"github.com/VladimirGrebenev/CarCare-backend/internal/usecase"
	"github.com/stretchr/testify/assert"
)

type mockUserRepo struct {
	users map[user.Email]*user.User
}

func (m *mockUserRepo) Create(ctx context.Context, u *user.User) error {
	if _, exists := m.users[u.Email]; exists {
		return user.ErrAlreadyExists
	}
	m.users[u.Email] = u
	return nil
}
func (m *mockUserRepo) GetByEmail(ctx context.Context, email user.Email) (*user.User, error) {
	u, ok := m.users[email]
	if !ok {
		return nil, user.ErrNotFound
	}
	return u, nil
}
func (m *mockUserRepo) GetByID(ctx context.Context, id string) (*user.User, error) {
	return nil, user.ErrNotFound
}
func (m *mockUserRepo) Update(ctx context.Context, u *user.User) error {
	return nil
}
func (m *mockUserRepo) Delete(ctx context.Context, id string) error {
	return nil
}
func (m *mockUserRepo) List(ctx context.Context) ([]*user.User, error) {
	return nil, nil
}

// ...другие методы user.Repository

func newAuthUsecase(repo *mockUserRepo) *usecase.AuthUsecase {
	return usecase.NewAuthUsecase(repo, &stubEmailSender{}, &stubJWT{}, &stubLogger{}, &stubSession{})
}

func TestRegister_Success(t *testing.T) {
	repo := &mockUserRepo{users: map[user.Email]*user.User{}}
	uc := newAuthUsecase(repo)
	token, err := uc.Register(context.Background(), "test@example.com", "password")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestRegister_Duplicate(t *testing.T) {
	repo := &mockUserRepo{users: map[user.Email]*user.User{user.Email("test@example.com"): {Email: user.Email("test@example.com"), Name: "test", Role: user.RoleUser}}}
	uc := newAuthUsecase(repo)
	_, err := uc.Register(context.Background(), "test@example.com", "password")
	assert.Error(t, err)
}

// ...тесты ConfirmEmail, Login, ForgotPassword, ResetPassword, OAuth, RefreshToken, Logout
