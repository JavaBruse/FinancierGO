package repositories

import (
	"database/sql"
	"financierGo/internal/models"
	"time"
)

type PaymentScheduleRepository struct {
	DB *sql.DB
}

func (r *PaymentScheduleRepository) Create(schedule *models.PaymentSchedule) error {
	query := `INSERT INTO payment_schedules (credit_id, amount, due_date, paid)
			  VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(query, schedule.CreditID, schedule.Amount, schedule.DueDate, schedule.Paid)
	return err
}

func (r *PaymentScheduleRepository) GetByCreditID(creditID int64) ([]models.PaymentSchedule, error) {
	rows, err := r.DB.Query(`SELECT id, credit_id, amount, due_date, paid FROM payment_schedules WHERE credit_id = $1`, creditID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.PaymentSchedule
	for rows.Next() {
		var s models.PaymentSchedule
		rows.Scan(&s.ID, &s.CreditID, &s.Amount, &s.DueDate, &s.Paid)
		schedules = append(schedules, s)
	}
	return schedules, nil
}

func (r *PaymentScheduleRepository) MarkPaid(id int64) error {
	_, err := r.DB.Exec(`UPDATE payment_schedules SET paid = true WHERE id = $1`, id)
	return err
}

func (r *PaymentScheduleRepository) GetAllDueTo(accountID int64, until time.Time) ([]models.PaymentSchedule, error) {
	query := `
		SELECT s.id, s.credit_id, s.amount, s.due_date, s.paid
		FROM payment_schedules s
		JOIN credits c ON s.credit_id = c.id
		WHERE c.account_id = $1 AND s.paid = false AND s.due_date <= $2`
	rows, err := r.DB.Query(query, accountID, until)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.PaymentSchedule
	for rows.Next() {
		var s models.PaymentSchedule
		rows.Scan(&s.ID, &s.CreditID, &s.Amount, &s.DueDate, &s.Paid)
		result = append(result, s)
	}
	return result, nil
}
