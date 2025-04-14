package repositories

import (
	"database/sql"
	"financierGo/internal/models"
	"time"
)

type TransactionRepository struct {
	DB *sql.DB
}

func (r *TransactionRepository) Create(t *models.Transaction) error {
	query := `INSERT INTO transactions (account_id, amount, type, date)
			  VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(query, t.AccountID, t.Amount, t.Type, t.Date)
	return err
}

func (r *TransactionRepository) GetForMonth(accountID int64, month time.Time) ([]models.Transaction, error) {
	start := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)
	rows, err := r.DB.Query(`SELECT id, account_id, amount, type, date FROM transactions 
							 WHERE account_id = $1 AND date >= $2 AND date < $3`,
		accountID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []models.Transaction
	for rows.Next() {
		var t models.Transaction
		rows.Scan(&t.ID, &t.AccountID, &t.Amount, &t.Type, &t.Date)
		txs = append(txs, t)
	}
	return txs, nil
}
