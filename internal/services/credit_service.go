package services

import (
	"errors"
	"math"
	"time"

	"financierGo/internal/models"
	"financierGo/internal/repositories"
)

type CreditService struct {
	Repo         *repositories.CreditRepository
	ScheduleRepo *repositories.PaymentScheduleRepository
	AccountRepo  *repositories.AccountRepository
}

func (s *CreditService) Create(accountID int64, amount, rate float64, months int) (*models.Credit, error) {
	account, err := s.AccountRepo.GetByID(accountID)
	if err != nil || account == nil {
		return nil, errors.New("account not found")
	}

	monthly := calcAnnuity(amount, rate, months)
	start := time.Now()
	credit := &models.Credit{
		AccountID:   accountID,
		Amount:      amount,
		Rate:        rate,
		Months:      months,
		StartDate:   start,
		NextPayment: start.Add(30 * 24 * time.Hour),
		Remaining:   amount,
	}

	if err := s.Repo.Create(credit); err != nil {
		return nil, err
	}

	// Генерация графика
	for i := 1; i <= months; i++ {
		p := &models.PaymentSchedule{
			CreditID: credit.ID,
			Amount:   monthly,
			DueDate:  start.Add(time.Duration(i) * 30 * 24 * time.Hour),
			Paid:     false,
		}
		s.ScheduleRepo.Create(p)
	}
	return credit, nil
}

func calcAnnuity(P, r float64, n int) float64 {
	monthlyRate := r / 100 / 12
	return P * monthlyRate / (1 - math.Pow(1+monthlyRate, -float64(n)))
}
