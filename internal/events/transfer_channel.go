package events

import (
	"goBank/internal/models"
	"goBank/internal/repository"
)

var (
	TransferCreateChannel   = make(chan models.Transfer)
	TransferResponseChannel = make(chan models.Transfer)
)

func PersistTransferWorker() {
	for transfer := range TransferCreateChannel {
		createdTransfer, _ := repository.SaveTransfer(transfer)
		TransferResponseChannel <- createdTransfer
	}
}
