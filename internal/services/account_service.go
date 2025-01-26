package services

import (
	"github.com/google/uuid"
	"goBank/internal/models"
	"goBank/internal/repository"
	"math/rand"
	"strconv"
)

func CreateAccount(account models.Account) (models.Account, error) {
	account.ID = uuid.New().String()
	account.Number = strconv.Itoa(10000000 + rand.Intn(99999999-10000000))
	account.SortCode = "001942"

	createdAccount, err := repository.SaveAccount(account)
	if err != nil {
		return models.Account{}, err
	}
	return createdAccount, nil
}

func GetAccountById(id string) (models.Account, error) {
	return repository.GetAccountById(id)
}

func GetAccounts() models.Accounts {
	return models.Accounts{Accounts: repository.GetAllAccounts()}
}
