package events

import (
	"goBank/internal/models"
	"goBank/internal/repository"
)

var (
	AccountCreateChannel   = make(chan models.Account)
	AccountResponseChannel = make(chan models.Account)
)

func PersistAccountWorker() {
	for account := range AccountCreateChannel {
		createdAccount, _ := repository.SaveAccount(account)
		AccountResponseChannel <- createdAccount
	}
}
