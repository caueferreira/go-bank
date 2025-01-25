package events

import (
	"goBank/internal/models"
	"goBank/internal/repository"
)

var (
	TransactionCreateChannel   = make(chan models.Transaction)
	TransactionResponseChannel = make(chan models.Transaction)
	TransactionErrorChannel    = make(chan error)
)

func PersistTransactionWorker() {
	for transaction := range TransactionCreateChannel {
		if transaction.TransactionType == "CREDIT" {
			_, err := repository.CreditAccount(transaction)
			if err != nil {
				TransactionErrorChannel <- err
				return
			}
		} else {
			_, err := repository.DebitAccount(transaction)
			if err != nil {
				TransactionErrorChannel <- err
				return
			}
		}
		createdTransaction, _ := repository.SaveTransaction(transaction)
		TransactionResponseChannel <- createdTransaction
	}
}
