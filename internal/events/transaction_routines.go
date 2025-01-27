package events

import (
	"goBank/internal/models"
	"goBank/internal/services"
)

var (
	GetAllTransactionsRoutine         = make(chan struct{})
	GetAllTransactionsResponseRoutine = make(chan models.Transactions)
	TransactionCreateRoutine          = make(chan models.CreateTransaction)
	TransactionResponseRoutine        = make(chan models.Transaction)
	TransactionErrorRoutine           = make(chan error)
)

func PersistTransactionWorker() {
	for transaction := range TransactionCreateRoutine {
		var createdTransaction models.Transaction
		var err error
		if transaction.TransactionType == "CREDIT" {
			createdTransaction, err = services.CreateCredit(transaction)
			if err != nil {
				TransactionErrorRoutine <- err
				return
			}
		} else {
			createdTransaction, err = services.CreateDebit(transaction)
			if err != nil {
				TransactionErrorRoutine <- err
				return
			}
		}
		TransactionResponseRoutine <- createdTransaction
	}
}

func GetTransactionsWorker() {
	for range GetAllTransactionsRoutine {
		GetAllTransactionsResponseRoutine <- services.GetTransactions()
	}
}
