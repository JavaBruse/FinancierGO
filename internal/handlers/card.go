package handlers

import (
	"encoding/json"
	"financierGo/internal/middleware"
	"financierGo/internal/services"
	"net/http"
)

type CardHandler struct {
	Service services.ICardService
}

func (h *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccountID int64  `json:"account_id"`
		Type      string `json:"type"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	userID := middleware.GetUserID(r)
	card, err := h.Service.CreateCard(userID, req.AccountID, req.Type)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(card)
}
