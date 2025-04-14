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
)

type mockCBRService struct {
	getKeyRateFunc func() (*models.KeyRate, error)
}

func (m *mockCBRService) GetKeyRate() (*models.KeyRate, error) {
	return m.getKeyRateFunc()
}

func TestGetKeyRateHandler(t *testing.T) {
	tests := []struct {
		name           string
		mockService    services.ICBRService
		expectedStatus int
	}{
		{
			name: "Successful key rate retrieval",
			mockService: &mockCBRService{
				getKeyRateFunc: func() (*models.KeyRate, error) {
					return &models.KeyRate{
						Rate: 16.0,
						Date: "2024-04-14",
					}, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Error getting key rate",
			mockService: &mockCBRService{
				getKeyRateFunc: func() (*models.KeyRate, error) {
					return nil, errors.New("failed to get key rate")
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/cbr/key-rate", nil)
			ctx := context.WithValue(req.Context(), middleware.UserKey, int64(1))
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler := &CBRHandler{
				Service: tt.mockService,
			}

			handler.GetKeyRate(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}
