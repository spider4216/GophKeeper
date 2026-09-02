package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	shrModel "github.com/spider4216/GophKeeper/internal/model"
	"github.com/spider4216/GophKeeper/internal/server/config"
	"github.com/spider4216/GophKeeper/internal/server/models"
	"github.com/spider4216/GophKeeper/internal/server/services"
)

const (
	syncGetLimit int = 1
)

type Handler struct {
	cfg     *config.Config
	service *services.Service
	logger  *slog.Logger
}

func New(cfg *config.Config, logger *slog.Logger, service *services.Service) Handler {
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
		h.logger.Error("failed read body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req shrModel.RegisterReq

	if err := json.Unmarshal(body, &req); err != nil {
		h.logger.Error("unmarshall error", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := h.service.CreateUser(ctx, req.Email, req.Pass)

	if err != nil && h.service.IsErrAsDuplicate(err) {
		h.logger.Debug("Duplicate user")
		w.WriteHeader(http.StatusConflict)
		return
	}

	if err != nil {
		h.logger.Error("cannot create user", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := shrModel.RegisterResp{
		UserID: userID,
	}

	b, err := json.Marshal(res)
	if err != nil {
		h.logger.Error("marshall error", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if _, err := w.Write(b); err != nil {
		h.logger.Error("failed to write response", "error", err)
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
		h.logger.Error("failed read body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req := shrModel.LoginReq{}

	if err := json.Unmarshal(body, &req); err != nil {
		h.logger.Error("failed read unmarshal", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Извлечение пользователя по логину
	user, err := h.service.GetUserByEmail(ctx, req.Email)
	if err != nil {
		h.logger.Error("User not found", "error", err)
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
		h.logger.Error("Unathorized: %s", "error", err)
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
		h.logger.Error("failed marshal: %s", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(b); err != nil {
		h.logger.Error("failed to write response: %s", "error", err)
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
		h.logger.Error("failed read body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req shrModel.SyncPutReq

	if err := json.Unmarshal(body, &req); err != nil {
		h.logger.Error("failed read unmarshal", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := h.service.GetUserIdFromCtx(ctx)

	if err := h.service.ApplySync(ctx, req.Changes, userID); err != nil {
		h.logger.Error("failed sync", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rev, err := h.service.GetLatestUserRev(ctx, userID)
	if err != nil {
		h.logger.Error("cannot get latest user revision", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := shrModel.SyncPutResp{
		LastRev: rev,
	}

	b, err := json.Marshal(resp)
	if err != nil {
		h.logger.Error("failed marshal", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(b); err != nil {
		h.logger.Error("failed to write response", "error", err)
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
		h.logger.Error("cannot convert since to int", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")

	limit := syncGetLimit

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			h.logger.Error("cannot convert limit to int", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	userID := h.service.GetUserIdFromCtx(ctx)

	resp, err := h.service.SyncGet(ctx, userID, since, limit)
	if err != nil {
		h.logger.Error("cannot sync", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(resp)
	if err != nil {
		h.logger.Error("cannot marshal response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(b); err != nil {
		h.logger.Error("failed to write response", "error", err)
	}
}

func (h Handler) SyncGetChunk(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	itemID := r.PathValue("itemID")
	chunkNumRaw := r.PathValue("num")

	if itemID == "" || chunkNumRaw == "" {
		h.logger.Error("cannot get itemID and chunk num")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// userID := h.service.GetUserIdFromCtx(ctx)

	// todo получить item по ID и пользователю, если нету - то ошибка

	chunkNum, err := strconv.Atoi(chunkNumRaw)
	if err != nil {
		h.logger.Error("cannot convert chunk num to int")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chunk, err := h.service.GetItemChunk(ctx, itemID, chunkNum)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.logger.Error("record not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		h.logger.Error("cannot get chunk", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := shrModel.ChunkGetResp{
		Ciphertext: chunk.Ciphertext,
		ChunkNum:   chunk.ChunkNumber,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		h.logger.Error("cannot marshal response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(data); err != nil {
		h.logger.Error("failed to write response", "error", err)
	}
}

func (h Handler) SaveChunk(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	r.Body = http.MaxBytesReader(w, r.Body, h.cfg.ChunkBodySize)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("failed read body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req models.ChunkPutReq

	if err := json.Unmarshal(body, &req); err != nil {
		h.logger.Error("failed read unmarshal", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	itemID := r.PathValue("itemID")
	chunkNumRaw := r.PathValue("num")

	if itemID == "" || chunkNumRaw == "" {
		h.logger.Error("cannot get itemID and chunk num")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// userID := h.service.GetUserIdFromCtx(ctx)

	// todo получить item по ID и пользователю, если нету - то ошибка

	chunkNum, err := strconv.Atoi(chunkNumRaw)
	if err != nil {
		h.logger.Error("cannot convert chunk num to int")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.service.InsertItemChunk(ctx, itemID, chunkNum, req.Ciphertext); err != nil {
		h.logger.Error("cannot insert item chunk")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if _, err := w.Write(nil); err != nil {
		h.logger.Error("failed to write response", "error", err)
	}
}
