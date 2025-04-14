package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"financierGo/internal/middleware"
	"financierGo/internal/services"

	"github.com/gorilla/mux"
)

type AnalyticsHandler struct {
	Service services.IAnalyticsService
}

func (h *AnalyticsHandler) Stats(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	stats, err := h.Service.GetStats(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stats)
}

func (h *AnalyticsHandler) CreditLoad(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	load, err := h.Service.GetCreditLoad(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(load)
}

func (h *AnalyticsHandler) Predict(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID, err := strconv.ParseInt(vars["accountId"], 10, 64)
	if err != nil {
		http.Error(w, "invalid account ID", http.StatusBadRequest)
		return
	}

	userID := middleware.GetUserID(r)
	prediction, err := h.Service.Predict(accountID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(prediction)
}
