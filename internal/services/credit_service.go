package services

import (
	"errors"
	"github.com/google/uuid"
	"goBank/internal/db"
	"goBank/internal/models"
)

func CreateCredit(credit models.Credit) (models.Account, error) {
	credit.ID = uuid.New().String()

	db.AccountsMutex.Lock()
	account, exists := db.Accounts[credit.AccountId]
	if !exists {
		return models.Account{}, errors.New("Account not found")
	}

	account.Balance += credit.Amount
	db.Accounts[account.ID] = account
	db.AccountsMutex.Unlock()

	return account, nil
}
