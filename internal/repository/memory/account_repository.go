package memory

import (
	"errors"
	"goBank/internal/db"
	"goBank/internal/models"
)

func SaveAccount(account models.Account) (models.Account, error) {
	db.AccountsMutex.Lock()
	db.Accounts[account.ID] = account
	db.AccountsMutex.Unlock()
	return account, nil
}

func GetAccountById(id string) (models.Account, error) {
	db.AccountsMutex.Lock()
	account, _ := db.Accounts[id]
	db.AccountsMutex.Unlock()

	return account, nil
}

func CreditAccount(credit models.Transaction) (models.Account, error) {
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
func DebitAccount(debit models.Transaction) (models.Account, error) {
	db.AccountsMutex.Lock()
	account, exists := db.Accounts[debit.AccountId]
	if !exists {
		return models.Account{}, errors.New("account not found")
	}
	if account.Balance < debit.Amount {
		return models.Account{}, errors.New("not enough balance")
	}
	account.Balance -= debit.Amount
	db.Accounts[account.ID] = account
	db.AccountsMutex.Unlock()
	return account, nil
}

func GetAllAccounts() []models.Account {
	db.AccountsMutex.Lock()
	var accounts []models.Account
	for _, account := range db.Accounts {
		accounts = append(accounts, account)
	}
	db.AccountsMutex.Unlock()
	return accounts
}
