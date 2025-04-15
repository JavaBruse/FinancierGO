package models

import "time"

type Credit struct {
	ID          int64     `json:"id"`
	AccountID   int64     `json:"account_id"`
	Amount      float64   `json:"amount"`
	Rate        float64   `json:"rate"`
	Months      int       `json:"months"`
	StartDate   time.Time `json:"start_date"`
	NextPayment time.Time `json:"next_payment"`
	Remaining   float64   `json:"remaining"`
}
