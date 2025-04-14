package handlers

import (
	"encoding/json"
	"financierGo/internal/services"
	"financierGo/internal/utils"
	"financierGo/pkg/logger"
	"net/http"
)

type AuthHandler struct {
	Service services.IUserService
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	// Логируем попытку регистрации
	logger.LogUserAction(0, "REGISTER", "Attempting to register user: "+req.Email)

	user, err := h.Service.Register(req.Username, req.Email, req.Password)
	if err != nil {
		logger.LogError(0, "REGISTER", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Логируем успешную регистрацию
	logger.LogDBWrite("users", user.ID, "New user registered: "+req.Email)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	// Логируем попытку входа
	logger.LogUserAction(0, "LOGIN", "Attempting to login: "+req.Email)

	user, err := h.Service.Login(req.Email, req.Password)
	if err != nil {
		logger.LogError(0, "LOGIN", "Unauthorized access attempt: "+req.Email)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Логируем успешный вход
	logger.LogUserAction(user.ID, "LOGIN", "Successfully logged in")
	token, _ := utils.GenerateJWT(user.ID)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
