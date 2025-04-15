package repositories

import (
	"database/sql"
	"financierGo/internal/models"
)

type AccountRepository struct {
	DB *sql.DB
}

func (r *AccountRepository) GetUserEmail(accountID int64) string {
	query := `SELECT u.email FROM users u 
			  JOIN accounts a ON a.user_id = u.id WHERE a.id = $1`
	var email string
	_ = r.DB.QueryRow(query, accountID).Scan(&email)
	return email
}

func (r *AccountRepository) Create(account *models.Account) error {
	query := `INSERT INTO accounts (user_id, number, balance, currency, created_at) 
			  VALUES ($1, $2, $3, $4, NOW()) RETURNING id`
	return r.DB.QueryRow(query, account.UserID, account.Number, account.Balance, account.Currency).
		Scan(&account.ID)
}

func (r *AccountRepository) GetByID(accountID int64) (*models.Account, error) {
	row := r.DB.QueryRow(`SELECT id, user_id, number, balance, currency, created_at 
						  FROM accounts WHERE id = $1`, accountID)

	var acc models.Account
	err := row.Scan(&acc.ID, &acc.UserID, &acc.Number, &acc.Balance, &acc.Currency, &acc.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &acc, err
}

func (r *AccountRepository) UpdateBalance(accountID int64, newBalance float64) error {
	_, err := r.DB.Exec(`UPDATE accounts SET balance = $1 WHERE id = $2`, newBalance, accountID)
	return err
}
