package services

import (
	"errors"
	"fmt"

	"financierGo/config"
	"financierGo/internal/models"
	"financierGo/internal/repositories"
	"financierGo/internal/utils"
)

type CardService struct {
	Repo        *repositories.CardRepository
	AccountRepo *repositories.AccountRepository
}

func (s *CardService) CreateCard(accountID, userID int64, cvv string) (*models.Card, error) {
	cfg := config.Load()

	account, err := s.AccountRepo.GetByID(accountID)
	if err != nil || account == nil || account.UserID != userID {
		return nil, errors.New("account not found or forbidden")
	}

	number := utils.GenerateCardNumber()
	expiry := utils.GenerateCardExpiry()
	cardData := fmt.Sprintf("%s|%s", number, expiry)

	encrypted, err := utils.EncryptPGP(cardData)
	if err != nil {
		return nil, err
	}

	cvvHash, err := utils.HashCVV(cvv)
	if err != nil {
		return nil, err
	}

	hmacSecret := cfg.JWT.HMAC
	hmac := utils.GenerateHMAC(cardData, hmacSecret)

	card := &models.Card{
		AccountID: accountID,
		Encrypted: encrypted,
		CVVHash:   cvvHash,
		HMAC:      hmac,
	}

	err = s.Repo.Create(card)
	return card, err
}
