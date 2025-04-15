package models

import "time"

type Card struct {
	ID        int64     `json:"id"`
	AccountID int64     `json:"account_id"`
	Encrypted string    `json:"encrypted"`
	CVVHash   string    `json:"-"`
	HMAC      string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}
