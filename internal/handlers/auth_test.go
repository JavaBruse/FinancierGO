package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"financierGo/internal/models"
	"financierGo/internal/services"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockUserService struct {
	createUserFunc func(username, email, password string) (*models.User, error)
	loginFunc      func(email, password string) (*models.User, error)
}

func (m *mockUserService) Register(username, email, password string) (*models.User, error) {
	return m.createUserFunc(username, email, password)
}

func (m *mockUserService) Login(email, password string) (*models.User, error) {
	return m.loginFunc(email, password)
}

func TestRegisterHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]string
		mockService    services.IUserService
		expectedStatus int
	}{
		{
			name: "Successful registration",
			requestBody: map[string]string{
				"username": "testuser",
				"email":    "test@example.com",
				"password": "password123",
			},
			mockService: &mockUserService{
				createUserFunc: func(username, email, password string) (*models.User, error) {
					return &models.User{
						ID:       1,
						Username: username,
						Email:    email,
					}, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid email format",
			requestBody: map[string]string{
				"username": "testuser",
				"email":    "invalid-email",
				"password": "password123",
			},
			mockService: &mockUserService{
				createUserFunc: func(username, email, password string) (*models.User, error) {
					return nil, errors.New("invalid email format")
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handler := &AuthHandler{
				Service: tt.mockService,
			}

			handler.Register(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}

func TestLoginHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]string
		mockService    services.IUserService
		expectedStatus int
	}{
		{
			name: "Successful login",
			requestBody: map[string]string{
				"email":    "test@example.com",
				"password": "password123",
			},
			mockService: &mockUserService{
				loginFunc: func(email, password string) (*models.User, error) {
					return &models.User{
						ID:       1,
						Username: "testuser",
						Email:    email,
					}, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid credentials",
			requestBody: map[string]string{
				"email":    "test@example.com",
				"password": "wrongpassword",
			},
			mockService: &mockUserService{
				loginFunc: func(email, password string) (*models.User, error) {
					return nil, errors.New("invalid credentials")
				},
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handler := &AuthHandler{
				Service: tt.mockService,
			}

			handler.Login(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var response map[string]string
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Errorf("error unmarshaling response: %v", err)
				}
				if _, exists := response["token"]; !exists {
					t.Error("response does not contain token")
				}
			}
		})
	}
}
