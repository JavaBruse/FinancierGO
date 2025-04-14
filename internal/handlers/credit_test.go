package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"financierGo/internal/middleware"
	"financierGo/internal/models"
	"financierGo/internal/services"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

type mockCreditService struct {
	createCreditFunc func(userID int64, amount float64, term int, rate float64) (*models.Credit, error)
	getScheduleFunc  func(creditID int64, userID int64) ([]models.Payment, error)
}

func (m *mockCreditService) CreateCredit(userID int64, amount float64, term int, rate float64) (*models.Credit, error) {
	if m.createCreditFunc != nil {
		return m.createCreditFunc(userID, amount, term, rate)
	}
	return nil, nil
}

func (m *mockCreditService) GetSchedule(creditID int64, userID int64) ([]models.Payment, error) {
	if m.getScheduleFunc != nil {
		return m.getScheduleFunc(creditID, userID)
	}
	return nil, nil
}

func TestCreateCreditHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		userID         int64
		mockService    services.ICreditService
		expectedStatus int
	}{
		{
			name: "Successful credit creation",
			requestBody: map[string]interface{}{
				"amount":        10000,
				"term":          12,
				"interest_rate": 10.5,
			},
			userID: 1,
			mockService: &mockCreditService{
				createCreditFunc: func(userID int64, amount float64, term int, rate float64) (*models.Credit, error) {
					return &models.Credit{
						ID:          1,
						AccountID:   1,
						Amount:      amount,
						Rate:        rate,
						Months:      term,
						StartDate:   time.Now(),
						NextPayment: time.Now().AddDate(0, 1, 0),
						Remaining:   amount,
					}, nil
				},
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid amount",
			requestBody: map[string]interface{}{
				"amount":        -1000,
				"term":          12,
				"interest_rate": 10.5,
			},
			userID: 1,
			mockService: &mockCreditService{
				createCreditFunc: func(userID int64, amount float64, term int, rate float64) (*models.Credit, error) {
					return nil, errors.New("invalid amount")
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/credits", bytes.NewBuffer(body))
			ctx := context.WithValue(req.Context(), middleware.UserKey, tt.userID)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler := NewCreditHandler(tt.mockService)
			handler.Create(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestGetScheduleHandler(t *testing.T) {
	tests := []struct {
		name           string
		creditID       string
		userID         int64
		mockService    services.ICreditService
		expectedStatus int
	}{
		{
			name:     "Successful schedule retrieval",
			creditID: "1",
			userID:   1,
			mockService: &mockCreditService{
				getScheduleFunc: func(creditID int64, userID int64) ([]models.Payment, error) {
					return []models.Payment{
						{
							ID:       1,
							CreditID: 1,
							Amount:   1000,
							Date:     time.Now().Format("2006-01-02"),
							Status:   "pending",
						},
					}, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "Credit not found",
			creditID: "999",
			userID:   1,
			mockService: &mockCreditService{
				getScheduleFunc: func(creditID int64, userID int64) ([]models.Payment, error) {
					return nil, errors.New("credit not found")
				},
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := mux.NewRouter()
			handler := NewCreditHandler(tt.mockService)
			router.HandleFunc("/api/credits/{creditId}/schedule", handler.GetSchedule).Methods("GET")

			req := httptest.NewRequest("GET", "/api/credits/"+tt.creditID+"/schedule", nil)
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
