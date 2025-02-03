package events

import (
	kafka2 "goBank/internal/events/kafka"
	"goBank/internal/models"
	"goBank/internal/services"
)

var (
	GetAllTransfersRoutine         = make(chan struct{})
	FindTransferRoutine            = make(chan string)
	GetAllTransfersResponseRoutine = make(chan models.Transfers)
	TransferCreateRoutine          = make(chan models.Transfer)
	TransferResponseRoutine        = make(chan models.Transfer)
)

func PersistTransferWorker() {
	for transfer := range TransferCreateRoutine {
		createdTransfer, _ := services.CreateTransfer(transfer)
		TransferResponseRoutine <- createdTransfer
	}
}

func FindTransferWorker() {
	for id := range FindTransferRoutine {
		transferFound, _ := services.GetTransferById(id)
		TransferResponseRoutine <- transferFound
	}
}

func GetAllTransferWorker() {
	for range GetAllTransfersRoutine {
		GetAllTransfersResponseRoutine <- services.GetTransfers()
	}
}

func KafkaTransferCreateWorker(consumer *kafka2.KafkaConsumer[models.Transfer]) {
	go func() {
		for msg := range consumer.Messages {
			go services.CreateTransfer(msg.Message)
		}
	}()
}
