package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spider4216/GophKeeper/internal/server/models"
	"github.com/spider4216/GophKeeper/internal/server/repositories"
	"go.uber.org/zap"
)

type Service struct {
	repo   repositories.Repository
	logger *zap.SugaredLogger
}

// todo moce to another file
type claims struct {
	jwt.RegisteredClaims
	UserID int64
}

func New(repo repositories.Repository, logger *zap.SugaredLogger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// todo раскидать методы по файлам

func (s *Service) CreateUser(ctx context.Context, email string, plainPass string) (int64, error) {
	// todo посолить пароль
	b := sha256.Sum256([]byte(plainPass))
	hash := hex.EncodeToString(b[:])

	u := models.UserRepo{
		Email:        email,
		PasswordHash: hash,
	}

	return s.repo.CreateUser(ctx, u)
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*models.UserRepo, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *Service) CheckPass(user *models.UserRepo, plainPass string) bool {
	hash := sha256.Sum256([]byte(plainPass))
	hashString := hex.EncodeToString(hash[:])

	return hashString == user.PasswordHash
}

func (s *Service) BuildJWTString(userId int64, secret string, exp time.Duration) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
		UserID: userId,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	s.logger.Debug("User ID is ", userId, " set to token")

	// возвращаем строку токена
	return tokenString, nil
}
