package services

import (
	"errors"
	"github.com/google/uuid"
	"goBank/internal/events"
	"goBank/internal/models"
	"goBank/internal/repository"
	"time"
)

func CreateDebit(newDebit models.Transaction) (models.Account, error) {
	if newDebit.TransactionType != "DEBIT" {
		return models.Account{}, errors.New("invalid transaction type")
	}

	newDebit.ID = uuid.New().String()

	events.TransactionCreateChannel <- newDebit

	select {
	case debit := <-events.TransactionResponseChannel:
		return repository.GetAccountById(debit.AccountId)
	case <-time.After(5 * time.Second):
		return models.Account{}, errors.New("persistence operation timed out")
	}
}
