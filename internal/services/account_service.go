package services

import (
	"errors"
	"github.com/google/uuid"
	"goBank/internal/db"
	"goBank/internal/models"
	"math/rand"
	"strconv"
)

func CreateAccount(account models.Account) (models.Account, error) {
	db.AccountsMutex.Lock()
	account.ID = uuid.New().String()
	account.Number = strconv.Itoa(10000000 + rand.Intn(99999999-10000000))
	account.SortCode = "001942"

	db.Accounts[account.ID] = account
	db.AccountsMutex.Unlock()
	return account, nil
}

func GetAccountById(id string) (models.Account, error) {
	db.AccountsMutex.Lock()
	account, exists := db.Accounts[id]
	db.AccountsMutex.Unlock()

	if !exists {
		return models.Account{}, errors.New("Account does not exist")
	}

	return account, nil
}

func GetAccounts() models.Accounts {
	db.AccountsMutex.Lock()
	var accountList []models.Account
	for _, account := range db.Accounts {
		accountList = append(accountList, account)
	}
	accounts := models.Accounts{}
	accounts.Accounts = accountList
	db.AccountsMutex.Unlock()
	return accounts
}
