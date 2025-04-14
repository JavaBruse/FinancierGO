package repositories

import (
	"database/sql"
	"financierGo/internal/models"
)

type CreditRepository struct {
	DB *sql.DB
}

func (r *CreditRepository) Create(credit *models.Credit) error {
	query := `INSERT INTO credits (account_id, amount, rate, months, start_date, next_payment, remaining)
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	return r.DB.QueryRow(query, credit.AccountID, credit.Amount, credit.Rate,
		credit.Months, credit.StartDate, credit.NextPayment, credit.Remaining).Scan(&credit.ID)
}

func (r *CreditRepository) GetAll() ([]*models.Credit, error) {
	rows, err := r.DB.Query(`SELECT id, account_id, amount, rate, months, start_date, next_payment, remaining FROM credits`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var credits []*models.Credit
	for rows.Next() {
		var c models.Credit
		rows.Scan(&c.ID, &c.AccountID, &c.Amount, &c.Rate, &c.Months, &c.StartDate, &c.NextPayment, &c.Remaining)
		credits = append(credits, &c)
	}
	return credits, nil
}
