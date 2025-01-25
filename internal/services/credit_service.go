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

	transaction := models.Transaction{}
	transaction.ID = credit.ID
	transaction.Amount = credit.Amount
	transaction.AccountId = credit.AccountId
	transaction.Type = "CREDIT"

	db.TransactionsMutex.Lock()
	db.Transactions[transaction.ID] = transaction
	db.TransactionsMutex.Unlock()

	return account, nil
}
