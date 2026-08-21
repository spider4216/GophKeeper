package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/spider4216/GophKeeper/internal/server/config"
	"github.com/spider4216/GophKeeper/internal/server/models"
	"github.com/spider4216/GophKeeper/internal/server/services"
	"go.uber.org/zap"
)

type Handler struct {
	cfg     *config.Config
	service *services.Service
	logger  *zap.SugaredLogger
}

func New(cfg *config.Config, logger *zap.SugaredLogger, service *services.Service) Handler {
	return Handler{
		cfg:     cfg,
		service: service,
		logger:  logger,
	}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	r.Body = http.MaxBytesReader(w, r.Body, h.cfg.MaxBodySize)

	body, err := io.ReadAll(r.Body)

	if err != nil {
		h.logger.Errorf("failed read body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req models.UserReq

	if err := json.Unmarshal(body, &req); err != nil {
		h.logger.Errorf("unmarshall error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// todo валидация почты и пароля
	// todo валидация отсутствия почты в БД (уникальность)

	userID, err := h.service.CreateUser(ctx, req.Email, req.Pass)

	if err != nil {
		h.logger.Errorf("cannot create user: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := models.RegisterResp{
		UserID: userID,
	}

	b, err := json.Marshal(res)

	if err != nil {
		h.logger.Errorf("marshall error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if _, err := w.Write(b); err != nil {
		h.logger.Errorf("failed to write response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Авторизация пользователя
// todo move to another file
func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Получаем тело запроса
	r.Body = http.MaxBytesReader(w, r.Body, h.cfg.MaxBodySize)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Errorf("failed read body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req := models.LoginReq{}

	if err := json.Unmarshal(body, &req); err != nil {
		h.logger.Errorf("failed read unmarshal: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Извлечение пользователя по логину
	user, err := h.service.GetUserByEmail(ctx, req.Email)
	if err != nil {
		h.logger.Errorf("User not found: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Проверка пароля
	if !h.service.CheckPass(user, req.Pass) {
		h.logger.Error("password is incorrect")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Авторизация
	token, err := h.service.BuildJWTString(user.ID, h.cfg.JWTKey, h.cfg.ExpToken)
	if err != nil {
		h.logger.Errorf("Unathorized: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	now := time.Now()

	res := models.LoginResp{
		// todo status to const
		Status:    "success",
		Token:     token,
		Message:   "Login successfully",
		ExpiredAt: now.Add(h.cfg.ExpToken).Unix(),
		CreatedAt: now.Unix(),
	}

	b, err := json.Marshal(res)

	if err != nil {
		h.logger.Errorf("failed marshal: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(b); err != nil {
		h.logger.Errorf("failed to write response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (h Handler) SyncIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Получаем тело запроса
	r.Body = http.MaxBytesReader(w, r.Body, h.cfg.MaxBodySize)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Errorf("failed read body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req models.SyncInReq

	if err := json.Unmarshal(body, &req); err != nil {
		h.logger.Errorf("failed read unmarshal: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := h.service.GetUserIdFromCtx(ctx)

	if err := h.service.CreateItems(ctx, req.Changes, userID); err != nil {
		h.logger.Errorf("failed sync: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rev, err := h.service.GetLatestUserRev(ctx, userID)

	if err != nil {
		h.logger.Errorf("cannot get latest user revision: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := models.SyncInResp{
		LastRev: rev,
	}

	b, err := json.Marshal(resp)

	if err != nil {
		h.logger.Errorf("failed marshal: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(b); err != nil {
		h.logger.Errorf("failed to write response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
