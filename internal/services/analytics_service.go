package services

import (
	"time"

	"financierGo/internal/repositories"
)

type AnalyticsService struct {
	TxRepo       *repositories.TransactionRepository
	CreditRepo   *repositories.CreditRepository
	ScheduleRepo *repositories.PaymentScheduleRepository
}

func (s *AnalyticsService) MonthlyStats(accountID int64) (income, expense float64, err error) {
	now := time.Now()
	txs, err := s.TxRepo.GetForMonth(accountID, now)
	if err != nil {
		return
	}
	for _, t := range txs {
		if t.Amount > 0 {
			income += t.Amount
		} else {
			expense += -t.Amount
		}
	}
	return
}

func (s *AnalyticsService) CreditLoad(accountID int64) (totalDebt float64, err error) {
	credits, err := s.CreditRepo.GetAll()
	if err != nil {
		return
	}
	for _, c := range credits {
		if c.AccountID == accountID {
			totalDebt += c.Remaining
		}
	}
	return
}

func (s *AnalyticsService) PredictBalance(accountID int64, days int) (float64, error) {
	future := time.Now().Add(time.Duration(days) * 24 * time.Hour)
	schedules, err := s.ScheduleRepo.GetAllDueTo(accountID, future)
	if err != nil {
		return 0, err
	}
	var total float64
	for _, s := range schedules {
		total += s.Amount
	}
	return total, nil // сумма платежей, которые будут списаны
}
