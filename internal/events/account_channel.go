package events

import (
	"goBank/internal/models"
	"goBank/internal/services"
)

var (
	GetAllAccountsChannel         = make(chan struct{})
	FindAccountChannel            = make(chan string)
	AccountCreateChannel          = make(chan models.Account)
	AccountResponseChannel        = make(chan models.Account)
	GetAllAccountsResponseChannel = make(chan models.Accounts)
)

func PersistAccountWorker() {
	for account := range AccountCreateChannel {
		createdAccount, _ := services.CreateAccount(account)
		AccountResponseChannel <- createdAccount
	}
}

func FindAccountWorker() {
	for id := range FindAccountChannel {
		accountFound, _ := services.GetAccountById(id)
		AccountResponseChannel <- accountFound
	}
}

func GetAllAccountsWorker() {
	for range GetAllAccountsChannel {
		GetAllAccountsResponseChannel <- services.GetAccounts()
	}
}
