package services

import (
	"errors"
	"github.com/google/uuid"
	"goBank/internal/events"
	"goBank/internal/models"
	"goBank/internal/repository"
	"time"
)

func CreateCredit(newCredit models.Transaction) (models.Account, error) {
	if newCredit.TransactionType != "CREDIT" {
		return models.Account{}, errors.New("invalid transaction type")
	}

	newCredit.ID = uuid.New().String()

	events.TransactionCreateChannel <- newCredit

	select {
	case credit := <-events.TransactionResponseChannel:
		return repository.GetAccountById(credit.AccountId)
	case <-time.After(5 * time.Second):
		return models.Account{}, errors.New("persistence operation timed out")
	}
}
