package services

import (
	"goBank/internal/models"
	"goBank/internal/repository/cassandra"
)

func GetTransactions() models.Transactions {
	return models.Transactions{Transactions: cassandra.GetAllTransactions()}
}

func CreateTransaction(transaction models.Transaction) (models.Transaction, error) {
	return cassandra.SaveTransaction(transaction)
}
