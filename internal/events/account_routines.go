package events

import (
	"goBank/internal/models"
	"goBank/internal/services"
)

var (
	GetAllAccountsRoutine         = make(chan struct{})
	FindAccountRoutine            = make(chan string)
	AccountCreateRoutine          = make(chan models.CreateAccount)
	AccountResponseRoutine        = make(chan models.Account)
	GetAllAccountsResponseRoutine = make(chan models.Accounts)
)

func PersistAccountWorker() {
	for account := range AccountCreateRoutine {
		createdAccount, _ := services.CreateAccount(account)
		AccountResponseRoutine <- createdAccount
	}
}

func FindAccountWorker() {
	for id := range FindAccountRoutine {
		accountFound, _ := services.GetAccountById(id)
		AccountResponseRoutine <- accountFound
	}
}

func GetAllAccountsWorker() {
	for range GetAllAccountsRoutine {
		GetAllAccountsResponseRoutine <- services.GetAccounts()
	}
}
