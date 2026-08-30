package services

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/spider4216/GophKeeper/internal/client/models"
	"github.com/spider4216/GophKeeper/internal/client/repositories"
	"github.com/spider4216/GophKeeper/internal/enum"
	shrModel "github.com/spider4216/GophKeeper/internal/model"
)

type Option func(*Service) error

type Service struct {
	client *http.Client
	host   string
	repo   repositories.Repository
	logger *slog.Logger
}

func WithHTTPClient(cli *http.Client) Option {
	return func(s *Service) error {
		s.client = cli
		return nil
	}
}

func WithHost(h string) Option {
	return func(s *Service) error {
		s.host = h
		return nil
	}
}

func WithRepo(r repositories.Repository) Option {
	return func(s *Service) error {
		s.repo = r
		return nil
	}
}

func WithLogger(l *slog.Logger) Option {
	return func(s *Service) error {
		s.logger = l
		return nil
	}
}

func New(opts ...Option) (*Service, error) {
	s := &Service{}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Service) CreateLoginPassItem(ctx context.Context, t enum.SecretType, data models.LoginPassReq, key string, userID int64) (string, error) {
	// Формат хранения для типа
	d := models.LoginPassFmt{
		Login: data.Login,
		Pass:  data.Pass,
	}

	b, err := json.Marshal(d)
	if err != nil {
		return "", err
	}

	encrypted, err := s.EncryptData(b, []byte(key))
	if err != nil {
		return "", err
	}

	item := models.ItemRepo{
		Type:       t,
		Ciphertext: encrypted,
		UserID:     userID,
	}

	return s.repo.CreateItem(ctx, item)
}

func (s *Service) EncryptData(data []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Создаем случайный одноразовый вектор (nonce)
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Шифруем данные и добавляем nonce в начало итогового байтового среза
	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return hex.EncodeToString(ciphertext), nil
}

func (s *Service) DecryptData(encryptedHex string, key []byte) ([]byte, error) {
	ciphertext, err := hex.DecodeString(encryptedHex)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("incorrect length")
	}

	// Извлекаем nonce и сами зашифрованные данные
	nonce, cipherTextData := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Расшифровываем
	plainText, err := gcm.Open(nil, nonce, cipherTextData, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

func (s *Service) CreateMeta(ctx context.Context, itemID string, k string, v string) (int64, error) {
	return s.repo.GetCommonRepo().CreateMeta(ctx, itemID, k, v)
}

func (s *Service) UpdateLastUserRev(ctx context.Context, userID int64, rev int64) error {
	return s.repo.UpdateLastUserRev(ctx, userID, rev)
}

func (s *Service) CreateLastUserRev(ctx context.Context, userID int64, rev int64) error {
	return s.repo.CreateLastUserRev(ctx, userID, rev)
}

func (s *Service) GetLatestUserRev(ctx context.Context, userID int64) (int64, error) {
	return s.repo.GetLatestUserRev(ctx, userID)
}

func (s *Service) GetUserItems(ctx context.Context, userID int64) ([]models.ItemRepo, error) {
	return s.repo.GetUserItems(ctx, userID)
}

func (s *Service) GetUserItemsWithMeta(ctx context.Context, userID int64) ([]models.ItemWithMeta, error) {
	items, err := s.GetUserItems(ctx, userID)
	if err != nil {
		return nil, err
	}

	var itemIDs []string

	for _, item := range items {
		itemIDs = append(itemIDs, item.ID)
	}

	meta, err := s.repo.GetCommonRepo().GetMetadataByItemIDs(ctx, itemIDs)
	if err != nil {
		return nil, fmt.Errorf("cannot get meta by item ids: %w", err)
	}

	return s.buildItemsWithMeta(items, meta), nil
}

func (s *Service) GetUserItemByID(ctx context.Context, itemID string, userID int64) (*models.ItemRepo, error) {
	return s.repo.GetUserItemByID(ctx, itemID, userID)
}

func (s *Service) DeleteUserItem(ctx context.Context, itemID string, userID int64) error {
	return s.repo.DeleteUserItem(ctx, itemID, userID)
}

func (s *Service) UpdateLoginPass(ctx context.Context, itemID string, userID int64, login string, pass string, key string, title string, metaID int64) error {
	data := models.LoginPassFmt{
		Login: login,
		Pass:  pass,
	}

	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("loginpass marshal err: %w", err)
	}

	encrypted, err := s.EncryptData(b, []byte(key))
	if err != nil {
		return fmt.Errorf("cannot enctypt data: %w", err)
	}

	if err := s.repo.UpdateLoginPass(ctx, itemID, userID, encrypted, metaID, title); err != nil {
		return fmt.Errorf("cannot update login and password: %w", err)
	}

	return nil
}

func (s *Service) GetMetadataByItemID(ctx context.Context, itemID string) ([]shrModel.MetadataRepo, error) {
	return s.repo.GetMetadataByItemID(ctx, itemID)
}

func (s *Service) CreateUserPassItem(ctx context.Context, data models.LoginPassReq, key string, userID int64) error {
	// Формат хранения для типа
	d := models.LoginPassFmt{
		Login: data.Login,
		Pass:  data.Pass,
	}

	b, err := json.Marshal(d)
	if err != nil {
		return err
	}

	encrypted, err := s.EncryptData(b, []byte(key))
	if err != nil {
		return err
	}

	item := models.ItemRepo{
		ID:         uuid.NewString(),
		Type:       enum.LoginPass,
		Ciphertext: encrypted,
		UserID:     userID,
	}

	return s.repo.CreateUserPassItem(ctx, item, userID, data.Title)
}
