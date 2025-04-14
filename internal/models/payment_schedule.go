package models

import "time"

type PaymentSchedule struct {
	ID       int64     `json:"id"`
	CreditID int64     `json:"credit_id"`
	Amount   float64   `json:"amount"`
	DueDate  time.Time `json:"due_date"`
	Paid     bool      `json:"paid"`
}
