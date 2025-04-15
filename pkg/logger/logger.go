package logger

import (
	"log"
)

// LogUserAction логирует действие пользователя
func LogUserAction(userID int64, action, details string) {
	log.Printf("[USER ACTION] UserID: %d | Action: %s | Details: %s",
		userID, action, details)
}

// LogDBWrite логирует запись в базу данных
func LogDBWrite(table string, userID int64, details string) {
	log.Printf("[DB WRITE] Table: %s | UserID: %d | Details: %s",
		table, userID, details)
}

// LogDBRead логирует чтение из базы данных
func LogDBRead(table string, userID int64, details string) {
	log.Printf("[DB READ] Table: %s | UserID: %d | Details: %s",
		table, userID, details)
}

// LogError логирует ошибку
func LogError(userID int64, action, errorMsg string) {
	log.Printf("[ERROR] UserID: %d | Action: %s | Error: %s",
		userID, action, errorMsg)
}
