package models

import "time"

type Transaction struct {
	ID        int64     `json:"id"`
	AccountID int64     `json:"account_id"`
	Amount    float64   `json:"amount"`
	Type      string    `json:"type"`
	Date      time.Time `json:"date"`
}
