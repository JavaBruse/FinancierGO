package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtKey []byte

func SetJWTSecret(secret string) {
	jwtKey = []byte(secret)
}

func GenerateJWT(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseJWT(tokenStr string) (int64, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return int64(claims["user_id"].(float64)), nil
	}
	return 0, err
}
