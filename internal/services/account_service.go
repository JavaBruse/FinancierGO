package services

import (
	"errors"
	"financierGo/internal/models"
	"financierGo/internal/repositories"
	"fmt"
	"math/rand"
	"time"
)

type AccountService struct {
	Repo *repositories.AccountRepository
}

func (s *AccountService) Create(userID int64, currency string) (*models.Account, error) {
	number := generateAccountNumber()
	account := &models.Account{
		UserID:   userID,
		Number:   number,
		Balance:  0.0,
		Currency: currency,
	}

	err := s.Repo.Create(account)
	return account, err
}

func (s *AccountService) Transfer(fromID, toID int64, amount float64, userID int64) error {
	from, err := s.Repo.GetByID(fromID)
	if err != nil || from == nil || from.UserID != userID {
		return errors.New("invalid source account")
	}
	if from.Balance < amount {
		return errors.New("insufficient funds")
	}

	to, err := s.Repo.GetByID(toID)
	if err != nil || to == nil {
		return errors.New("invalid destination account")
	}

	from.Balance -= amount
	to.Balance += amount

	err = s.Repo.UpdateBalance(from.ID, from.Balance)
	if err != nil {
		return err
	}
	return s.Repo.UpdateBalance(to.ID, to.Balance)
}

func generateAccountNumber() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("40817%010d", rand.Int63n(1e10))
}
