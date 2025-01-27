package events

import (
	"goBank/internal/models"
	"goBank/internal/services"
)

var (
	GetAllTransactionsChannel         = make(chan struct{})
	GetAllTransactionsResponseChannel = make(chan models.Transactions)
	TransactionCreateChannel          = make(chan models.CreateTransaction)
	TransactionResponseChannel        = make(chan models.Transaction)
	TransactionErrorChannel           = make(chan error)
)

func PersistTransactionWorker() {
	for transaction := range TransactionCreateChannel {
		var createdTransaction models.Transaction
		var err error
		if transaction.TransactionType == "CREDIT" {
			createdTransaction, err = services.CreateCredit(transaction)
			if err != nil {
				TransactionErrorChannel <- err
				return
			}
		} else {
			createdTransaction, err = services.CreateDebit(transaction)
			if err != nil {
				TransactionErrorChannel <- err
				return
			}
		}
		//createdTransaction, _ := services.CreateTransaction(transaction)
		TransactionResponseChannel <- createdTransaction
	}
}

func GetTransactionsWorker() {
	for range GetAllTransactionsChannel {
		GetAllTransactionsResponseChannel <- services.GetTransactions()
	}
}
