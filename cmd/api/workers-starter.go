package main

import (
	"goBank/internal/events"
	"goBank/internal/events/kafka"
)

func StartWorkers() {
	go events.PersistAccountWorker()
	go events.FindAccountWorker()
	go events.GetAllAccountsWorker()

	kafka.InitTransferKafkaHandlers()
	//go events.PersistTransferWorker()
	go events.FindTransferWorker()
	go events.GetAllTransferWorker()

	go events.PersistTransactionWorker()
	go events.GetTransactionsWorker()
}
