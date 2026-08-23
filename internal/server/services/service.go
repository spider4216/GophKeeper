package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v4"
	shrModel "github.com/spider4216/GophKeeper/internal/model"
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

// todo move to other file
type (
	userIdKey string
)

// todo move to other file
const (
	userKey userIdKey = "user_id"
)

func New(repo repositories.Repository, logger *zap.SugaredLogger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// todo move to other file
func (s *Service) SetUserIdToCtx(ctx context.Context, userId int64) context.Context {
	return context.WithValue(ctx, userKey, userId)
}

// todo move to other file
func (s *Service) GetUserIdFromCtx(ctx context.Context) int64 {
	userId, ok := ctx.Value(userKey).(int64)

	if !ok {
		return 0
	}

	return userId
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

func (s *Service) CreateItems(ctx context.Context, in []shrModel.SyncPutChange, userID int64) error {
	// todo сделать гоурутины
	for _, item := range in {

		s.logger.Debugf("Sync Operation: %s", item.Operation)

		line := models.ItemRepo{
			ID:     int64(item.Item.ID),
			UserID: userID,
			Type:   item.Item.Type,
		}

		// todo transaction
		// todo в зависимости от типа операции разные действия
		// todo разделить на методы
		// todo op co conts
		if item.Operation == "CREATE" {
			s.logger.Debug("create strategy")
			itemID, err := s.repo.CreateItem(ctx, line)

			if err != nil {
				return err
			}

			pl := models.ItemPayloadRepo{
				ItemID:     itemID,
				Ciphertext: item.Item.Ciphertext,
			}

			if err := s.repo.CreateItemPayload(ctx, pl); err != nil {
				return err
			}

			for k, v := range item.Metadata {
				if _, err := s.repo.CreateMeta(ctx, itemID, k, v); err != nil {
					return err
				}
			}

			if _, err := s.repo.CreateSyncChanges(ctx, itemID, item.Operation, userID); err != nil {
				return err
			}
		}

		if item.Operation == "DELETE" {
			s.logger.Debug("Delete strategy")
			// todo transaction
			if err := s.repo.GetCommonRepo().DeleteUserMetaByItemID(ctx, int64(item.Item.ID), userID); err != nil {
				return err
			}

			if err := s.repo.DeletePayloadByItemID(ctx, int64(item.Item.ID)); err != nil {
				return err
			}

			if err := s.repo.GetCommonRepo().DeleteUserItemByID(ctx, int64(item.Item.ID), userID); err != nil {
				return err
			}

			// todo op to const
			if _, err := s.repo.CreateSyncChanges(ctx, int64(item.Item.ID), "DELETE", userID); err != nil {
				return err
			}
		}

		if item.Operation == "UPDATE" {
			// todo transaction
			s.logger.Debug("Update strategy")
			if err := s.repo.UpdateUserItemPayload(ctx, int64(item.Item.ID), userID, item.Item.Ciphertext); err != nil {
				return err
			}

			// todo update metadata

			if _, err := s.repo.CreateSyncChanges(ctx, int64(item.Item.ID), "UPDATE", userID); err != nil {
				return err
			}
		}

	}

	return nil
}

func (s *Service) GetLatestUserRev(ctx context.Context, userID int64) (int64, error) {
	return s.repo.GetLatestUserRev(ctx, userID)
}

func (s *Service) SyncGet(ctx context.Context, userID int64, since int64) (*shrModel.SyncGet, error) {
	s.logger.Debug("Get changes...")
	changes, err := s.repo.GetUserSyncChanges(ctx, userID, since)

	if err != nil {
		return nil, err
	}

	var itemIDs []int64

	for _, item := range changes {
		itemIDs = append(itemIDs, item.ItemID)
	}

	s.logger.Debug("Get Items...")

	items, err := s.repo.GetItemsByIDs(ctx, itemIDs)

	if err != nil {
		return nil, err
	}

	s.logger.Debug("Get payloads...")

	payloads, err := s.repo.GetPayloadByItemIDs(ctx, itemIDs)

	if err != nil {
		return nil, err
	}

	s.logger.Debug("Get metadata...")

	meta, err := s.repo.GetCommonRepo().GetMetadataByItemIDs(ctx, itemIDs)

	if err != nil {
		return nil, err
	}

	s.logger.Debug("Get latest rev...")

	rev, err := s.repo.GetLatestUserRev(ctx, userID)

	if err != nil {
		return nil, err
	}

	return s.mapSyncResponse(changes, items, payloads, meta, rev), nil

}
