package handlers

import (
	"encoding/json"
	"financierGo/internal/services"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type CreditHandler struct {
	Service *services.CreditService
}

func (h *CreditHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccountID int64   `json:"account_id"`
		Amount    float64 `json:"amount"`
		Rate      float64 `json:"rate"`
		Months    int     `json:"months"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	credit, err := h.Service.Create(req.AccountID, req.Amount, req.Rate, req.Months)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(credit)
}

func (h *CreditHandler) GetSchedule(w http.ResponseWriter, r *http.Request) {
	creditID, _ := strconv.ParseInt(mux.Vars(r)["creditId"], 10, 64)
	schedules, err := h.Service.ScheduleRepo.GetByCreditID(creditID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(schedules)
}
