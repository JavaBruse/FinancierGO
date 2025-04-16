package services

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"financierGo/internal/models"
	"financierGo/internal/repositories"
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

	tx, err := s.Repo.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	from.Balance -= amount
	to.Balance += amount

	_, err = tx.Exec(`UPDATE accounts SET balance = $1 WHERE id = $2`, from.Balance, from.ID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`UPDATE accounts SET balance = $1 WHERE id = $2`, to.Balance, to.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func generateAccountNumber() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("40817%010d", rand.Int63n(1e10))
}

func (s *AccountService) AdjustBalance(accountID, userID int64, delta float64) error {
	account, err := s.Repo.GetByID(accountID)
	if err != nil || account == nil || account.UserID != userID {
		return errors.New("нет доступа к счету")
	}

	if delta < 0 && account.Balance < -delta {
		return errors.New("недостаточно средств")
	}

	account.Balance += delta
	return s.Repo.UpdateBalance(account.ID, account.Balance)
}
