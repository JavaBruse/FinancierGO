package models

import "time"

type Card struct {
	ID        int64     `json:"id"`
	AccountID int64     `json:"account_id"`
	Encrypted string    `json:"encrypted"` // номер + срок действия в PGP
	CVVHash   string    `json:"-"`         // хэш CVV
	HMAC      string    `json:"-"`         // контрольная сумма
	CreatedAt time.Time `json:"created_at"`
}
