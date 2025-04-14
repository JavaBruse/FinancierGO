package models

import "time"

type Transaction struct {
	ID        int64     `json:"id"`
	AccountID int64     `json:"account_id"`
	Amount    float64   `json:"amount"` // положительное — доход, отрицательное — расход
	Type      string    `json:"type"`   // например: "income", "expense", "credit_payment"
	Date      time.Time `json:"date"`
}
