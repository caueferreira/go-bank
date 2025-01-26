package services

import (
	"errors"
	"github.com/google/uuid"
	"goBank/internal/models"
	"goBank/internal/repository/cassandra"
	"time"
)

func CreateDebit(newDebit models.Transaction) (models.Account, error) {
	if newDebit.TransactionType != "DEBIT" {
		return models.Account{}, errors.New("invalid transaction type")
	}

	newDebit.ID = uuid.New().String()
	newDebit.CreatedAt = time.Now().Unix()

	account, err := cassandra.DebitAccount(newDebit)
	if err != nil {
		return models.Account{}, err
	}
	cassandra.SaveTransaction(newDebit)
	return account, nil
}
