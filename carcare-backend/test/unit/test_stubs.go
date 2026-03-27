
package unit

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/user"
)

type stubUserRepo struct{}

func (s *stubUserRepo) Create(ctx context.Context, u *user.User) error { return nil }
func (s *stubUserRepo) GetByID(ctx context.Context, id string) (*user.User, error) {
	return nil, nil
}
func (s *stubUserRepo) GetByEmail(ctx context.Context, email user.Email) (*user.User, error) {
	return nil, nil
}
func (s *stubUserRepo) Update(ctx context.Context, u *user.User) error  { return nil }
func (s *stubUserRepo) Delete(ctx context.Context, id string) error      { return nil }
func (s *stubUserRepo) List(ctx context.Context) ([]*user.User, error)   { return nil, nil }

// stubUserRepoWithUser — stub с заранее созданным пользователем для тестов Login
type stubUserRepoWithUser struct {
	email        string
	passwordHash string
}

func newStubUserRepoWithUser(email, plainPassword string) *stubUserRepoWithUser {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.MinCost)
	return &stubUserRepoWithUser{email: email, passwordHash: string(hashed)}
}

func (s *stubUserRepoWithUser) Create(ctx context.Context, u *user.User) error { return nil }
func (s *stubUserRepoWithUser) GetByID(ctx context.Context, id string) (*user.User, error) {
	return nil, errors.New("not found")
}
func (s *stubUserRepoWithUser) GetByEmail(ctx context.Context, email user.Email) (*user.User, error) {
	if string(email) == s.email {
		return &user.User{
			ID:           "stub-id",
			Email:        email,
			Name:         "Stub User",
			Role:         user.RoleUser,
			PasswordHash: s.passwordHash,
		}, nil
	}
	return nil, errors.New("not found")
}
func (s *stubUserRepoWithUser) Update(ctx context.Context, u *user.User) error { return nil }
func (s *stubUserRepoWithUser) Delete(ctx context.Context, id string) error    { return nil }
func (s *stubUserRepoWithUser) List(ctx context.Context) ([]*user.User, error) { return nil, nil }

type stubLogger struct{}

func (s *stubLogger) Info(args ...interface{})  {}
func (s *stubLogger) Error(args ...interface{}) {}

type stubEmailSender struct{}

func (s *stubEmailSender) SendConfirmation(email, token string) error  { return nil }
func (s *stubEmailSender) SendPasswordReset(email, token string) error { return nil }

type stubJWT struct{}

func (s *stubJWT) Generate(userID string) (string, error) { return "", nil }
func (s *stubJWT) Validate(token string) (string, error)  { return "", nil }

type stubSession struct{}

func (s *stubSession) Create(userID string) (string, error)        { return "", nil }
func (s *stubSession) Refresh(refreshToken string) (string, error) { return "", nil }
func (s *stubSession) Invalidate(token string) error               { return nil }
