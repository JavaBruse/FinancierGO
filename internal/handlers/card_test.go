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
)

type mockCardService struct {
	createCardFunc func(userID int64, accountID int64, cardType string) (*models.Card, error)
}

func (m *mockCardService) CreateCard(userID int64, accountID int64, cardType string) (*models.Card, error) {
	return m.createCardFunc(userID, accountID, cardType)
}

func TestCreateCardHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		userID         int64
		mockService    services.ICardService
		expectedStatus int
	}{
		{
			name: "Successful card creation",
			requestBody: map[string]interface{}{
				"account_id": 1,
				"type":       "debit",
			},
			userID: 1,
			mockService: &mockCardService{
				createCardFunc: func(userID int64, accountID int64, cardType string) (*models.Card, error) {
					return &models.Card{
						ID:        1,
						AccountID: accountID,
						Encrypted: "encrypted_card_data",
						CVVHash:   "hashed_cvv",
						HMAC:      "control_sum",
						CreatedAt: time.Now(),
					}, nil
				},
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid card type",
			requestBody: map[string]interface{}{
				"account_id": 1,
				"type":       "invalid",
			},
			userID: 1,
			mockService: &mockCardService{
				createCardFunc: func(userID int64, accountID int64, cardType string) (*models.Card, error) {
					return nil, errors.New("invalid card type")
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/cards", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			ctx := context.WithValue(req.Context(), middleware.UserKey, tt.userID)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler := &CardHandler{
				Service: tt.mockService,
			}

			handler.CreateCard(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}
