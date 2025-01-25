package services

import (
	"errors"
	"github.com/google/uuid"
	"goBank/internal/db"
	"goBank/internal/models"
)

func CreateTransfer(transfer models.Transfer) (models.Accounts, error) {
	transfer.ID = uuid.New().String()
	transfer.Success = false

	db.AccountsMutex.Lock()
	toAccount, toExists := db.Accounts[transfer.ToAccount]
	fromAccount, fromExists := db.Accounts[transfer.FromAccount]

	if !toExists || !fromExists {
		db.AccountsMutex.Unlock()
		return models.Accounts{}, errors.New("Account not found")
	}

	if fromAccount.Balance >= transfer.Amount {
		toAccount.Balance += transfer.Amount
		fromAccount.Balance -= transfer.Amount
		db.Accounts[toAccount.ID] = toAccount
		db.Accounts[fromAccount.ID] = fromAccount
		transfer.Success = true
	}
	accounts := models.Accounts{Accounts: []models.Account{toAccount, fromAccount}}
	db.AccountsMutex.Unlock()

	db.TransfersMutex.Lock()
	db.Transfers[transfer.ID] = transfer
	db.TransfersMutex.Unlock()

	return accounts, nil
}

func GetTransferById(id string) (models.Transfer, error) {
	db.TransfersMutex.Lock()
	transfer, exists := db.Transfers[id]
	db.TransfersMutex.Unlock()

	if !exists {
		return models.Transfer{}, errors.New("Transfer does not exist")
	}

	return transfer, nil
}

func GetTransfers() models.Transfers {
	db.TransfersMutex.Lock()
	var transferList []models.Transfer
	for _, transfer := range db.Transfers {
		transferList = append(transferList, transfer)
	}
	transfers := models.Transfers{}
	transfers.Transfers = transferList
	db.TransfersMutex.Unlock()

	return transfers
}
