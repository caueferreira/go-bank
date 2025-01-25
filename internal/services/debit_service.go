package services

import (
	"errors"
	"github.com/google/uuid"
	"goBank/internal/db"
	"goBank/internal/models"
)

func CreateDebit(debit models.Debit) (models.Account, error) {
	debit.ID = uuid.New().String()

	db.AccountsMutex.Lock()
	account, exists := db.Accounts[debit.AccountId]
	if !exists {
		return models.Account{}, errors.New("Account not found")
	}

	account.Balance -= debit.Amount
	db.Accounts[account.ID] = account
	db.AccountsMutex.Unlock()

	return account, nil
}
