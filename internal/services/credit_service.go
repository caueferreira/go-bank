package services

import (
	"errors"
	"github.com/google/uuid"
	"goBank/internal/models"
	"goBank/internal/repository"
)

func CreateCredit(newCredit models.Transaction) (models.Account, error) {
	if newCredit.TransactionType != "DEBIT" {
		return models.Account{}, errors.New("invalid transaction type")
	}

	newCredit.ID = uuid.New().String()

	account, err := repository.CreditAccount(newCredit)
	if err != nil {
		return models.Account{}, err
	}
	repository.SaveTransaction(newCredit)
	return account, nil
}
