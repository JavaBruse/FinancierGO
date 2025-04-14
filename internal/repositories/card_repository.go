package repositories

import (
	"database/sql"
	"financierGo/internal/models"
)

type CardRepository struct {
	DB *sql.DB
}

func (r *CardRepository) Create(card *models.Card) error {
	query := `INSERT INTO cards (account_id, encrypted, cvv_hash, hmac, created_at)
			  VALUES ($1, $2, $3, $4, NOW()) RETURNING id`
	return r.DB.QueryRow(query, card.AccountID, card.Encrypted, card.CVVHash, card.HMAC).
		Scan(&card.ID)
}

func (r *CardRepository) GetByID(id int64) (*models.Card, error) {
	row := r.DB.QueryRow(`SELECT id, account_id, encrypted, cvv_hash, hmac, created_at FROM cards WHERE id = $1`, id)

	var c models.Card
	err := row.Scan(&c.ID, &c.AccountID, &c.Encrypted, &c.CVVHash, &c.HMAC, &c.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &c, err
}
