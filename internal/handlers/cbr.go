package handlers

import (
	"encoding/json"
	"net/http"

	"financierGo/internal/services"
)

type CBRHandler struct {
	Service *services.CBRService
}

func (h *CBRHandler) GetKeyRate(w http.ResponseWriter, r *http.Request) {
	rate, err := h.Service.GetKeyRate()
	if err != nil {
		http.Error(w, "Ошибка получения ставки ЦБ: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]float64{"key_rate": rate})
}
