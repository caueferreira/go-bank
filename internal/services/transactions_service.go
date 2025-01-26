package services

import (
	"goBank/internal/models"
	"goBank/internal/repository"
)

func GetTransactions() models.Transactions {
	return models.Transactions{Transactions: repository.GetAllTransactions()}
}

func CreateTransaction(transaction models.Transaction) (models.Transaction, error) {
	return repository.SaveTransaction(transaction)
}
