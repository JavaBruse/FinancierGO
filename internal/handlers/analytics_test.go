package handlers

import (
	"context"
	"errors"
	"financierGo/internal/middleware"
	"financierGo/internal/models"
	"financierGo/internal/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

type mockAnalyticsService struct {
	getStatsFunc      func(userID int64) (*models.Stats, error)
	getCreditLoadFunc func(userID int64) (*models.CreditLoad, error)
	predictFunc       func(accountID int64, userID int64) (*models.Prediction, error)
}

func (m *mockAnalyticsService) GetStats(userID int64) (*models.Stats, error) {
	return m.getStatsFunc(userID)
}

func (m *mockAnalyticsService) GetCreditLoad(userID int64) (*models.CreditLoad, error) {
	return m.getCreditLoadFunc(userID)
}

func (m *mockAnalyticsService) Predict(accountID int64, userID int64) (*models.Prediction, error) {
	return m.predictFunc(accountID, userID)
}

func TestStatsHandler(t *testing.T) {
	tests := []struct {
		name           string
		userID         int64
		mockService    services.IAnalyticsService
		expectedStatus int
	}{
		{
			name:   "Successful stats retrieval",
			userID: 1,
			mockService: &mockAnalyticsService{
				getStatsFunc: func(userID int64) (*models.Stats, error) {
					return &models.Stats{
						TotalBalance:    100000.0,
						MonthlyIncome:   50000.0,
						MonthlyExpenses: 30000.0,
						ActiveCredits:   2,
					}, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Error getting stats",
			userID: 1,
			mockService: &mockAnalyticsService{
				getStatsFunc: func(userID int64) (*models.Stats, error) {
					return nil, errors.New("failed to get stats")
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/analytics", nil)
			ctx := context.WithValue(req.Context(), middleware.UserKey, tt.userID)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler := &AnalyticsHandler{
				Service: tt.mockService,
			}

			handler.Stats(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}

func TestCreditLoadHandler(t *testing.T) {
	tests := []struct {
		name           string
		userID         int64
		mockService    services.IAnalyticsService
		expectedStatus int
	}{
		{
			name:   "Successful credit load retrieval",
			userID: 1,
			mockService: &mockAnalyticsService{
				getCreditLoadFunc: func(userID int64) (*models.CreditLoad, error) {
					return &models.CreditLoad{
						TotalDebt:      200000.0,
						MonthlyPayment: 20000.0,
						DebtToIncome:   0.4,
					}, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Error getting credit load",
			userID: 1,
			mockService: &mockAnalyticsService{
				getCreditLoadFunc: func(userID int64) (*models.CreditLoad, error) {
					return nil, errors.New("failed to get credit load")
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/analytics/credit", nil)
			ctx := context.WithValue(req.Context(), middleware.UserKey, tt.userID)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler := &AnalyticsHandler{
				Service: tt.mockService,
			}

			handler.CreditLoad(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}

func TestPredictHandler(t *testing.T) {
	tests := []struct {
		name           string
		accountID      string
		userID         int64
		mockService    services.IAnalyticsService
		expectedStatus int
	}{
		{
			name:      "Successful prediction",
			accountID: "1",
			userID:    1,
			mockService: &mockAnalyticsService{
				predictFunc: func(accountID int64, userID int64) (*models.Prediction, error) {
					return &models.Prediction{
						NextMonthBalance: 120000.0,
						Trend:            "up",
						Confidence:       0.85,
					}, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "Account not found",
			accountID: "999",
			userID:    1,
			mockService: &mockAnalyticsService{
				predictFunc: func(accountID int64, userID int64) (*models.Prediction, error) {
					return nil, errors.New("account not found")
				},
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := mux.NewRouter()
			handler := &AnalyticsHandler{Service: tt.mockService}
			router.HandleFunc("/api/accounts/{accountId}/predict", handler.Predict).Methods("GET")

			req := httptest.NewRequest("GET", "/api/accounts/"+tt.accountID+"/predict", nil)
			ctx := context.WithValue(req.Context(), middleware.UserKey, tt.userID)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}
