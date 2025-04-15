package scheduler

import (
	"financierGo/internal/utils"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"financierGo/internal/repositories"
)

type CreditScheduler struct {
	CreditRepo   *repositories.CreditRepository
	ScheduleRepo *repositories.PaymentScheduleRepository
	AccountRepo  *repositories.AccountRepository
}

func (s *CreditScheduler) Start(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			s.processPayments()
		}
	}()
}

func (s *CreditScheduler) processPayments() {
	logrus.Infof("[Scheduler] Checking for due payments...")

	credits, _ := s.CreditRepo.GetAll()
	now := time.Now()

	for _, credit := range credits {
		schedules, _ := s.ScheduleRepo.GetByCreditID(credit.ID)

		for _, sched := range schedules {
			if !sched.Paid && sched.DueDate.Before(now) {
				account, _ := s.AccountRepo.GetByID(credit.AccountID)

				if account.Balance >= sched.Amount {
					// списание
					account.Balance -= sched.Amount
					credit.Remaining -= sched.Amount
					s.AccountRepo.UpdateBalance(account.ID, account.Balance)
					s.ScheduleRepo.MarkPaid(sched.ID)
					logrus.Warnf("Списано %.2f по кредиту %d", sched.Amount, credit.ID)
				} else {
					// штраф
					penalty := sched.Amount * 0.1
					credit.Remaining += penalty

					// отправка письма
					userEmail := s.AccountRepo.GetUserEmail(account.ID)
					body := fmt.Sprintf("Здравствуйте! На вашем счете недостаточно средств для оплаты кредита. Начислен штраф +10%%. Сумма: %.2f руб.", penalty)
					err := utils.SendEmail(userEmail, "Просрочка по кредиту", body)
					if err != nil {
						logrus.WithError(err).Error("Ошибка отправки email")
					}
				}
			}
		}
	}
}
