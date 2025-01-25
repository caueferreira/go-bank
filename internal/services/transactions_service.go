package services

import (
	"goBank/internal/models"
	"goBank/internal/repository"
)

func GetTransactions() models.Transactions {
	return models.Transactions{Transactions: repository.GetAllTransactions()}
}
