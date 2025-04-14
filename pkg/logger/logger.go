package logger

import (
	"log"
	"time"
)

// LogUserAction логирует действие пользователя
func LogUserAction(userID int64, action, details string) {
	log.Printf("[USER ACTION] UserID: %d | Action: %s | Details: %s | Time: %s",
		userID, action, details, time.Now().Format("2006-01-02 15:04:05"))
}

// LogDBWrite логирует запись в базу данных
func LogDBWrite(table string, userID int64, details string) {
	log.Printf("[DB WRITE] Table: %s | UserID: %d | Details: %s | Time: %s",
		table, userID, details, time.Now().Format("2006-01-02 15:04:05"))
}

// LogDBRead логирует чтение из базы данных
func LogDBRead(table string, userID int64, details string) {
	log.Printf("[DB READ] Table: %s | UserID: %d | Details: %s | Time: %s",
		table, userID, details, time.Now().Format("2006-01-02 15:04:05"))
}

// LogError логирует ошибку
func LogError(userID int64, action, errorMsg string) {
	log.Printf("[ERROR] UserID: %d | Action: %s | Error: %s | Time: %s",
		userID, action, errorMsg, time.Now().Format("2006-01-02 15:04:05"))
}
