package usecase

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

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

func NewAuthUsecase(userRepo user.Repository, emailSender EmailSender, jwtSvc JWTService, logger Logger, session SessionManager) *AuthUsecase {
	return &AuthUsecase{
		UserRepo:    userRepo,
		EmailSender: emailSender,
		JWT:         jwtSvc,
		Logger:      logger,
		Session:     session,
	}
}

// getJWTSecret читает JWT_SECRET из окружения, fallback на дефолтное значение для разработки
func getJWTSecret() string {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return s
	}
	return "carcare-dev-secret-change-in-prod"
}

// generateJWT создаёт подписанный JWT-токен для пользователя
func generateJWT(userID, email string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(getJWTSecret()))
}

func (uc *AuthUsecase) logInfo(args ...interface{}) {
	if uc.Logger != nil {
		uc.Logger.Info(args...)
	}
}

func (uc *AuthUsecase) logError(args ...interface{}) {
	if uc.Logger != nil {
		uc.Logger.Error(args...)
	}
}

// Register создаёт нового пользователя и возвращает JWT-токен
func (uc *AuthUsecase) Register(ctx context.Context, email, password string) (string, error) {
	uc.logInfo("Register called", email)

	// Проверка на существование пользователя с таким email
	existing, _ := uc.UserRepo.GetByEmail(ctx, user.Email(email))
	if existing != nil {
		return "", user.ErrAlreadyExists
	}

	newUser := &user.User{
		ID:           uuid.New().String(),
		Email:        user.Email(email),
		Name:         email, // имя по умолчанию = email, можно переопределить позднее
		Role:         user.RoleUser,
		PasswordHash: password, // Create() сделает bcrypt
	}

	if err := uc.UserRepo.Create(ctx, newUser); err != nil {
		uc.logError("Register: failed to create user", err)
		return "", err
	}

	tokenStr, err := generateJWT(newUser.ID, email)
	if err != nil {
		uc.logError("Register: failed to generate JWT", err)
		return "", err
	}

	return tokenStr, nil
}

// Login проверяет credentials и возвращает JWT-токен
func (uc *AuthUsecase) Login(ctx context.Context, email, password string) (string, error) {
	uc.logInfo("Login called", email)

	u, err := uc.UserRepo.GetByEmail(ctx, user.Email(email))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	tokenStr, err := generateJWT(u.ID, string(u.Email))
	if err != nil {
		uc.logError("Login: failed to generate JWT", err)
		return "", err
	}

	return tokenStr, nil
}

func (uc *AuthUsecase) ConfirmEmail(ctx context.Context, token string) error {
	uc.logInfo("ConfirmEmail called", token)
	// TODO: найти пользователя по токену, активировать, удалить токен
	return nil
}

func (uc *AuthUsecase) ResendConfirmation(ctx context.Context, email string) error {
	uc.logInfo("ResendConfirmation called", email)
	// TODO: найти пользователя, сгенерировать новый токен, отправить email
	return nil
}

func (uc *AuthUsecase) ForgotPassword(ctx context.Context, email string) error {
	uc.logInfo("ForgotPassword called", email)
	// TODO: найти пользователя, сгенерировать токен, отправить email
	return nil
}

func (uc *AuthUsecase) ResetPassword(ctx context.Context, token, newPassword string) error {
	uc.logInfo("ResetPassword called", token)
	// TODO: найти пользователя по токену, обновить пароль, удалить токен
	return nil
}

func (uc *AuthUsecase) OAuthYandex(ctx context.Context, code string) (string, error) {
	uc.logInfo("OAuthYandex called", code)
	return "", nil
}

func (uc *AuthUsecase) OAuthGoogle(ctx context.Context, code string) (string, error) {
	uc.logInfo("OAuthGoogle called", code)
	return "", nil
}

func (uc *AuthUsecase) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	uc.logInfo("RefreshToken called", refreshToken)
	// TODO: валидация refreshToken, генерация нового accessToken
	return "", nil
}

func (uc *AuthUsecase) Logout(ctx context.Context, token string) error {
	uc.logInfo("Logout called", token)
	// TODO: инвалидировать refresh/session
	return nil
}
