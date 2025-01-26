package events

import (
	"goBank/internal/models"
	"goBank/internal/services"
)

var (
	GetAllTransfersChannel         = make(chan struct{})
	FindTransferChannel            = make(chan string)
	GetAllTransfersResponseChannel = make(chan models.Transfers)
	TransferCreateChannel          = make(chan models.Transfer)
	TransferResponseChannel        = make(chan models.Transfer)
)

func PersistTransferWorker() {
	for transfer := range TransferCreateChannel {
		createdTransfer, _ := services.CreateTransfer(transfer)
		TransferResponseChannel <- createdTransfer
	}
}

func FindTransferWorker() {
	for id := range FindTransferChannel {
		transferFound, _ := services.GetTransferById(id)
		TransferResponseChannel <- transferFound
	}
}

func GetAllTransferWorker() {
	for range GetAllTransfersChannel {
		GetAllTransfersResponseChannel <- services.GetTransfers()
	}
}
