package models

import "time"

type Account struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Number    string    `json:"number"` // Уникальный номер счёта
	Balance   float64   `json:"balance"`
	Currency  string    `json:"currency"` // Например: "RUB"
	CreatedAt time.Time `json:"created_at"`
}
