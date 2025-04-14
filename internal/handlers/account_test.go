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
)

type mockAccountService struct {
	createAccountFunc func(userID int64, currency string) (*models.Account, error)
	transferFunc      func(fromID, toID int64, amount float64, userID int64) error
}

func (m *mockAccountService) Create(userID int64, currency string) (*models.Account, error) {
	return m.createAccountFunc(userID, currency)
}

func (m *mockAccountService) Transfer(fromID, toID int64, amount float64, userID int64) error {
	return m.transferFunc(fromID, toID, amount, userID)
}

func TestCreateAccountHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		userID         int64
		mockService    services.IAccountService
		expectedStatus int
	}{
		{
			name: "Successful account creation",
			requestBody: map[string]interface{}{
				"currency": "RUB",
			},
			userID: 1,
			mockService: &mockAccountService{
				createAccountFunc: func(userID int64, currency string) (*models.Account, error) {
					return &models.Account{
						ID:       1,
						UserID:   userID,
						Number:   "4081700000000001",
						Balance:  0.0,
						Currency: currency,
					}, nil
				},
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid currency",
			requestBody: map[string]interface{}{
				"currency": "INVALID",
			},
			userID: 1,
			mockService: &mockAccountService{
				createAccountFunc: func(userID int64, currency string) (*models.Account, error) {
					return nil, errors.New("invalid currency")
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/accounts", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			ctx := context.WithValue(req.Context(), middleware.UserKey, tt.userID)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler := &AccountHandler{
				Service: tt.mockService,
			}

			handler.Create(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}

func TestTransferHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		userID         int64
		mockService    services.IAccountService
		expectedStatus int
	}{
		{
			name: "Successful transfer",
			requestBody: map[string]interface{}{
				"from_account_id": 1,
				"to_account_id":   2,
				"amount":          100.0,
			},
			userID: 1,
			mockService: &mockAccountService{
				transferFunc: func(fromID, toID int64, amount float64, userID int64) error {
					return nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Insufficient funds",
			requestBody: map[string]interface{}{
				"from_account_id": 1,
				"to_account_id":   2,
				"amount":          1000.0,
			},
			userID: 1,
			mockService: &mockAccountService{
				transferFunc: func(fromID, toID int64, amount float64, userID int64) error {
					return errors.New("insufficient funds")
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/transfer", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			ctx := context.WithValue(req.Context(), middleware.UserKey, tt.userID)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler := &AccountHandler{
				Service: tt.mockService,
			}

			handler.Transfer(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}
