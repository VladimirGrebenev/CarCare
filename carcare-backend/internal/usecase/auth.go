package usecase

import (
	"context"
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/user"
)

type AuthUsecase struct {
	UserRepo    user.Repository
	EmailSender EmailSender
	JWT         JWTService
	Logger      Logger
	Session     SessionManager
}

// EmailSender — интерфейс для отправки email (подтверждение, восстановление)
type EmailSender interface {
	SendConfirmation(email, token string) error
	SendPasswordReset(email, token string) error
}

// JWTService — интерфейс для генерации/валидации JWT
type JWTService interface {
	Generate(userID string) (string, error)
	Validate(token string) (string, error)
}

// Logger — интерфейс централизованного логирования
type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

// SessionManager — интерфейс для управления refresh/session
type SessionManager interface {
	Create(userID string) (string, error)
	Refresh(refreshToken string) (string, error)
	Invalidate(token string) error
}

func NewAuthUsecase(userRepo user.Repository, emailSender EmailSender, jwt JWTService, logger Logger, session SessionManager) *AuthUsecase {
	return &AuthUsecase{
		 UserRepo:    userRepo,
		 EmailSender: emailSender,
		 JWT:         jwt,
		 Logger:      logger,
		 Session:     session,
	}
}

func (uc *AuthUsecase) Register(ctx context.Context, email, password string) error {
	uc.Logger.Info("Register called", email)
	// ...
	return nil
}

func (uc *AuthUsecase) ConfirmEmail(ctx context.Context, token string) error {
	uc.Logger.Info("ConfirmEmail called", token)
	// TODO: найти пользователя по токену, активировать, удалить токен
	return nil
}

func (uc *AuthUsecase) ResendConfirmation(ctx context.Context, email string) error {
	uc.Logger.Info("ResendConfirmation called", email)
	// TODO: найти пользователя, сгенерировать новый токен, отправить email
	return nil
}

func (uc *AuthUsecase) Login(ctx context.Context, email, password string) (string, error) {
	uc.Logger.Info("Login called", email)
	// ...
	return "", nil
}

func (uc *AuthUsecase) ForgotPassword(ctx context.Context, email string) error {
	uc.Logger.Info("ForgotPassword called", email)
	// TODO: найти пользователя, сгенерировать токен, отправить email
	return nil
}

func (uc *AuthUsecase) ResetPassword(ctx context.Context, token, newPassword string) error {
	uc.Logger.Info("ResetPassword called", token)
	// TODO: найти пользователя по токену, обновить пароль, удалить токен
	return nil
}

func (uc *AuthUsecase) OAuthYandex(ctx context.Context, code string) (string, error) {
	uc.Logger.Info("OAuthYandex called", code)
	// ...
	return "", nil
}

func (uc *AuthUsecase) OAuthGoogle(ctx context.Context, code string) (string, error) {
	uc.Logger.Info("OAuthGoogle called", code)
	// ...
	return "", nil
}

func (uc *AuthUsecase) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	uc.Logger.Info("RefreshToken called", refreshToken)
	// TODO: валидация refreshToken, генерация нового accessToken
	return "", nil
}

func (uc *AuthUsecase) Logout(ctx context.Context, token string) error {
	uc.Logger.Info("Logout called", token)
	// TODO: инвалидировать refresh/session
	return nil
}
