package repository

import (
	"errors"
	"goBank/internal/db"
	"goBank/internal/models"
)

func SaveTransaction(transaction models.Transaction) (models.Transaction, error) {
	db.TransactionsMutex.Lock()
	db.Transactions[transaction.ID] = transaction
	db.TransactionsMutex.Unlock()
	return transaction, nil
}

func FindTransactionById(id string) (models.Transaction, error) {
	db.TransactionsMutex.Lock()
	transaction, exists := db.Transactions[id]
	db.TransactionsMutex.Unlock()

	if !exists {
		return models.Transaction{}, errors.New("transaction does not exist")
	}

	return transaction, nil
}

func GetAllTransactions() []models.Transaction {
	db.TransactionsMutex.Lock()
	var transactions []models.Transaction
	for _, transaction := range db.Transactions {
		transactions = append(transactions, transaction)
	}
	db.TransactionsMutex.Unlock()
	return transactions
}
