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

	newDebit := models.Debit{}
	newDebit.AccountId = transfer.FromAccount
	newDebit.Amount = transfer.Amount

	_, debitErr := CreateDebit(newDebit)

	newCredit := models.Credit{}
	newCredit.AccountId = transfer.ToAccount
	newCredit.Amount = transfer.Amount

	_, creditErr := CreateCredit(newCredit)

	if creditErr != nil {
		refund := models.Credit{}
		refund.AccountId = transfer.FromAccount
		refund.Amount = transfer.Amount

		_, err := CreateCredit(refund)
		if err != nil {
			return models.Accounts{}, err
		}
	}

	if debitErr != nil && creditErr != nil {
		transfer.Success = true
	}

	fromAccount, err := GetAccountById(transfer.FromAccount)
	if err != nil {
		return models.Accounts{}, err
	}
	toAccount, err := GetAccountById(transfer.ToAccount)
	if err != nil {
		return models.Accounts{}, err
	}

	accounts := models.Accounts{Accounts: []models.Account{toAccount, fromAccount}}

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
