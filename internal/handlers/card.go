package handlers

import (
	"encoding/json"
	"financierGo/internal/middleware"
	"financierGo/internal/services"
	"net/http"
)

type CardHandler struct {
	Service *services.CardService
}

func (h *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccountID int64  `json:"account_id"`
		CVV       string `json:"cvv"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	userID := middleware.GetUserID(r)
	card, err := h.Service.CreateCard(req.AccountID, userID, req.CVV)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]int64{"card_id": card.ID})
}
