package api

import (
	"context"
	"encoding/json"
	"net/http"

	"msg-processor/internal/entity"
)

type Service interface {
	CreateMsg(ctx context.Context, text string) (entity.Msg, error)
}

type Handler struct {
	s Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		s: s,
	}
}

type CreateMsgRequest struct {
	Text string `json:"text"`
}

func (h *Handler) CreateMsg(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateMsgRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	msg, err := h.s.CreateMsg(ctx, req.Text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
