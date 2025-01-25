package main

import "goBank/internal/events"

func StartWorkers() {
	go events.PersistAccountWorker()
	go events.PersistTransferWorker()
	go events.PersistTransactionWorker()
}
