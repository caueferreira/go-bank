package services

import (
	"errors"
	"github.com/google/uuid"
	"goBank/internal/events"
	"goBank/internal/models"
	"goBank/internal/repository"
	"math/rand"
	"strconv"
	"time"
)

func CreateAccount(account models.Account) (models.Account, error) {
	account.ID = uuid.New().String()
	account.Number = strconv.Itoa(10000000 + rand.Intn(99999999-10000000))
	account.SortCode = "001942"

	events.AccountCreateChannel <- account

	select {
	case account := <-events.AccountResponseChannel:
		return account, nil
	case <-time.After(5 * time.Second):
		return models.Account{}, errors.New("persistence operation timed out")
	}
}

func GetAccountById(id string) (models.Account, error) {
	return repository.GetAccountById(id)
}

func GetAccounts() models.Accounts {
	return models.Accounts{Accounts: repository.GetAllAccounts()}
}
