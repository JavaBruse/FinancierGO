package handlers

import (
	"encoding/json"
	"financierGo/internal/services"
	"financierGo/internal/utils"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	Service *services.UserService
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.WithError(err).Warn("Ошибка разбора тела запроса")
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	user, err := h.Service.Register(req.Username, req.Email, req.Password)
	if err != nil {
		logrus.WithError(err).Error("Ошибка при регистрации")

		if strings.Contains(err.Error(), "pq:") {
			http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	user, err := h.Service.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	token, _ := utils.GenerateJWT(user.ID)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
