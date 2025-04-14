package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"financierGo/internal/middleware"
	"financierGo/internal/services"
)

type AnalyticsHandler struct {
	Service *services.AnalyticsService
}

func (h *AnalyticsHandler) Stats(w http.ResponseWriter, r *http.Request) {
	accountID := middleware.GetUserID(r) // либо из запроса, как path param
	income, expense, err := h.Service.MonthlyStats(accountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]float64{
		"income":  income,
		"expense": expense,
	})
}

func (h *AnalyticsHandler) CreditLoad(w http.ResponseWriter, r *http.Request) {
	accountID := middleware.GetUserID(r)
	debt, err := h.Service.CreditLoad(accountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]float64{"debt": debt})
}

func (h *AnalyticsHandler) Predict(w http.ResponseWriter, r *http.Request) {
	accountID := middleware.GetUserID(r)
	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	amount, err := h.Service.PredictBalance(accountID, days)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]float64{"planned_expense": amount})
}
