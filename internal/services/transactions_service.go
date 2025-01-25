package services

import (
	"goBank/internal/db"
	"goBank/internal/models"
)

func GetTransactions() models.Transactions {
	db.TransactionsMutex.Lock()
	var transactions []models.Transaction
	for _, account := range db.Transactions {
		transactions = append(transactions, account)
	}
	response := models.Transactions{}
	response.Transactions = transactions
	db.TransactionsMutex.Unlock()
	return response
}
