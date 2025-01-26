package main

import "goBank/internal/events"

func StartWorkers() {
	go events.PersistAccountWorker()
	go events.FindAccountWorker()
	go events.GetAllAccountsWorker()

	go events.PersistTransferWorker()
	go events.FindTransferWorker()
	go events.GetAllTransferWorker()

	go events.PersistTransactionWorker()
	go events.GetTransactionsWorker()
}
