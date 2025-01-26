package services

import (
	"errors"
	"github.com/google/uuid"
	"goBank/internal/models"
	"goBank/internal/repository"
	"time"
)

func CreateCredit(newCredit models.Transaction) (models.Account, error) {
	if newCredit.TransactionType != "CREDIT" {
		return models.Account{}, errors.New("invalid transaction type")
	}

	newCredit.ID = uuid.New().String()
	newCredit.CreatedAt = time.Now().Unix()

	account, err := repository.CreditAccount(newCredit)
	if err != nil {
		return models.Account{}, err
	}
	repository.SaveTransaction(newCredit)
	return account, nil
}
