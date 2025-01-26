package services

import (
	"errors"
	"github.com/google/uuid"
	"goBank/internal/models"
	"goBank/internal/repository/cassandra"
	"time"
)

func CreateCredit(newCredit models.Transaction) (models.Account, error) {
	if newCredit.TransactionType != "CREDIT" {
		return models.Account{}, errors.New("invalid transaction type")
	}

	newCredit.ID = uuid.New().String()
	newCredit.CreatedAt = time.Now().Unix()

	account, err := cassandra.CreditAccount(newCredit)
	if err != nil {
		return models.Account{}, err
	}
	cassandra.SaveTransaction(newCredit)
	return account, nil
}
