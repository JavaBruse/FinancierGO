package handlers

import (
	"encoding/json"
	"financierGo/internal/middleware"
	"financierGo/internal/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CreditHandler struct {
	Service services.ICreditService
}

func NewCreditHandler(service services.ICreditService) *CreditHandler {
	return &CreditHandler{Service: service}
}

func (h *CreditHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Amount float64 `json:"amount"`
		Term   int     `json:"term"`
		Rate   float64 `json:"interest_rate"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := middleware.GetUserID(r)
	credit, err := h.Service.CreateCredit(userID, req.Amount, req.Term, req.Rate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(credit)
}

func (h *CreditHandler) GetSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	creditID, err := strconv.ParseInt(vars["creditId"], 10, 64)
	if err != nil {
		http.Error(w, "invalid credit ID", http.StatusBadRequest)
		return
	}

	userID := middleware.GetUserID(r)
	schedules, err := h.Service.GetSchedule(creditID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(schedules)
}
