package services

import (
	"errors"
	"github.com/google/uuid"
	"goBank/internal/models"
	"goBank/internal/repository/cassandra"
	"time"
)

func CreateCredit(newCredit models.CreateTransaction) (models.Transaction, error) {
	if newCredit.TransactionType != "CREDIT" {
		return models.Transaction{}, errors.New("invalid transaction type")
	}

	credit := models.Transaction{AccountId: newCredit.AccountId, Amount: newCredit.Amount}
	credit.ID = uuid.New().String()
	credit.CreatedAt = time.Now().Unix()

	_, err := cassandra.CreditAccount(credit)
	if err != nil {
		return models.Transaction{}, err
	}

	return cassandra.SaveTransaction(credit)
}
