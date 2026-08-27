package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	shrModel "github.com/spider4216/GophKeeper/internal/model"
	"github.com/spider4216/GophKeeper/internal/server/config"
	"github.com/spider4216/GophKeeper/internal/server/services"
	"go.uber.org/zap"
)

const (
	syncGetLimit int = 1
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

	var req shrModel.RegisterReq

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

	res := shrModel.RegisterResp{
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

	req := shrModel.LoginReq{}

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

	res := shrModel.LoginResp{
		Status:    "success",
		Token:     token,
		Message:   "Login successfully",
		ExpiredAt: now.Add(h.cfg.ExpToken).Unix(),
		CreatedAt: now.Unix(),
		UserID:    user.ID,
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

	var req shrModel.SyncPutReq

	if err := json.Unmarshal(body, &req); err != nil {
		h.logger.Errorf("failed read unmarshal: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := h.service.GetUserIdFromCtx(ctx)

	if err := h.service.ApplySync(ctx, req.Changes, userID); err != nil {
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

	resp := shrModel.SyncPutResp{
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

func (h Handler) SyncOut(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sinceStr := r.URL.Query().Get("since")
	if sinceStr == "" {
		h.logger.Error("no since found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	since, err := strconv.ParseInt(sinceStr, 10, 64)
	if err != nil {
		h.logger.Errorf("cannot convert since to int: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")

	limit := syncGetLimit

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			h.logger.Errorf("cannot convert limit to int: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	userID := h.service.GetUserIdFromCtx(ctx)

	resp, err := h.service.SyncGet(ctx, userID, since, limit)
	if err != nil {
		h.logger.Errorf("cannot sync: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(resp)
	if err != nil {
		h.logger.Errorf("cannot marshal response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(b); err != nil {
		h.logger.Errorf("failed to write response: %s", err)
	}
}
