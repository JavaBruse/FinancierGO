package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"financierGo/internal/middleware"
	"financierGo/internal/services"

	"github.com/gorilla/mux"
)

type AccountHandler struct {
	Service *services.AccountService
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Currency string `json:"currency"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	userID := middleware.GetUserID(r)
	account, err := h.Service.Create(userID, req.Currency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(account)
}

func (h *AccountHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	accountID, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	userID := middleware.GetUserID(r)

	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Amount <= 0 {
		http.Error(w, "Неверная сумма", http.StatusBadRequest)
		return
	}

	err := h.Service.AdjustBalance(accountID, userID, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *AccountHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	accountID, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	userID := middleware.GetUserID(r)

	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Amount <= 0 {
		http.Error(w, "Неверная сумма", http.StatusBadRequest)
		return
	}

	err := h.Service.AdjustBalance(accountID, userID, -req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *AccountHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FromAccountID int64   `json:"from_account_id"`
		ToAccountID   int64   `json:"to_account_id"`
		Amount        float64 `json:"amount"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	userID := middleware.GetUserID(r)
	err := h.Service.Transfer(req.FromAccountID, req.ToAccountID, req.Amount, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
