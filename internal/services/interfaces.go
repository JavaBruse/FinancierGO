package services

import "financierGo/internal/models"

type IAccountService interface {
	Create(userID int64, currency string) (*models.Account, error)
	Transfer(fromID, toID int64, amount float64, userID int64) error
}

type IUserService interface {
	Register(username, email, password string) (*models.User, error)
	Login(email, password string) (*models.User, error)
}

type ICardService interface {
	CreateCard(userID int64, accountID int64, cardType string) (*models.Card, error)
}

type ICreditService interface {
	CreateCredit(userID int64, amount float64, term int, rate float64) (*models.Credit, error)
	GetSchedule(creditID int64, userID int64) ([]models.Payment, error)
}

type IAnalyticsService interface {
	GetStats(userID int64) (*models.Stats, error)
	GetCreditLoad(userID int64) (*models.CreditLoad, error)
	Predict(accountID int64, userID int64) (*models.Prediction, error)
}

type ICBRService interface {
	GetKeyRate() (*models.KeyRate, error)
}
