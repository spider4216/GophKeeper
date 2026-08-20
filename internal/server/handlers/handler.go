package handlers

import (
	"encoding/json"
	"io"
	"net/http"

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

	if _, err := h.service.CreateUser(ctx, req.Email, req.Pass); err != nil {
		h.logger.Errorf("cannot create user: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if _, err := w.Write(nil); err != nil {
		h.logger.Errorf("failed to write response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
