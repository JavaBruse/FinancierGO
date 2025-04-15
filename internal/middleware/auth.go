package middleware

import (
	"context"
	"financierGo/internal/constants"
	"financierGo/internal/utils"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := utils.ParseJWT(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// Создаем новый контекст с user_id
		ctx := context.WithValue(r.Context(), constants.UserIDKey, userID)
		// Передаем запрос с новым контекстом
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(r *http.Request) int64 {
	return r.Context().Value(constants.UserIDKey).(int64)
}
