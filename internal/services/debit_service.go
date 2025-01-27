package services

import (
	"errors"
	"github.com/google/uuid"
	"goBank/internal/models"
	"goBank/internal/repository/cassandra"
	"time"
)

func CreateDebit(newDebit models.CreateTransaction) (models.Transaction, error) {
	if newDebit.TransactionType != "DEBIT" {
		return models.Transaction{}, errors.New("invalid transaction type")
	}

	debit := models.Transaction{AccountId: newDebit.AccountId, Amount: newDebit.Amount}
	debit.ID = uuid.New().String()
	debit.CreatedAt = time.Now().Unix()

	_, err := cassandra.DebitAccount(debit)
	if err != nil {
		return models.Transaction{}, err
	}

	return cassandra.SaveTransaction(debit)
}
