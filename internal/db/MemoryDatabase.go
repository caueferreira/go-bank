package db

import "sync"
import "goBank/internal/models"

var (
	Accounts      = make(map[string]models.Account)
	AccountsMutex = &sync.Mutex{}
)

var (
	Transfers      = make(map[string]models.Transfer)
	TransfersMutex = &sync.Mutex{}
)

var (
	Transactions      = make(map[string]models.Transaction)
	TransactionsMutex = &sync.Mutex{}
)
