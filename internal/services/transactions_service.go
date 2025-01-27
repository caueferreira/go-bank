package services

import (
	"github.com/google/uuid"
	"goBank/internal/models"
	"goBank/internal/repository/cassandra"
	"time"
)

func GetTransactions() models.Transactions {
	return models.Transactions{Transactions: cassandra.GetAllTransactions()}
}

func CreateTransaction(newTransaction models.CreateTransaction) (models.Transaction, error) {
	transaction := models.Transaction{AccountId: newTransaction.AccountId, Amount: newTransaction.Amount}
	transaction.ID = uuid.New().String()
	transaction.CreatedAt = time.Now().Unix()

	return cassandra.SaveTransaction(transaction)
}
