package services

import (
	"errors"
	"github.com/google/uuid"
	"goBank/internal/events"
	"goBank/internal/models"
	"goBank/internal/repository"
	"time"
)

func CreateTransfer(transfer models.Transfer) (models.Accounts, error) {
	transfer.ID = uuid.New().String()
	transfer.Success = false

	newDebit, newCredit, err := createTransactions(transfer)
	if err != nil {
		saveTransfer(transfer)
		return models.Accounts{}, err
	}

	_, debitErr := CreateDebit(newDebit)
	if debitErr != nil {
		saveTransfer(transfer)
		return models.Accounts{}, debitErr
	}

	_, creditErr := CreateCredit(newCredit)

	if creditErr != nil {
		refund := models.Transaction{}
		refund.AccountId = transfer.FromAccount
		refund.Amount = transfer.Amount
		refund.TransactionType = "CREDIT"

		_, err := CreateCredit(refund)
		if err != nil {
			saveTransfer(transfer)
			return models.Accounts{}, err
		}
	}

	if debitErr == nil && creditErr == nil {
		transfer.Success = true
	}

	saveTransfer(transfer)

	fromAccount, err := GetAccountById(transfer.FromAccount)
	if err != nil {
		return models.Accounts{}, err
	}
	toAccount, err := GetAccountById(transfer.ToAccount)
	if err != nil {
		return models.Accounts{}, err
	}

	accounts := models.Accounts{Accounts: []models.Account{toAccount, fromAccount}}

	return accounts, nil
}

func GetTransferById(id string) (models.Transfer, error) {
	return repository.GetTransferById(id)
}

func GetTransfers() models.Transfers {
	return models.Transfers{Transfers: repository.GetAllTransfers()}
}

func saveTransfer(transfer models.Transfer) (models.Transfer, error) {
	events.TransferCreateChannel <- transfer

	select {
	case transfer := <-events.TransferResponseChannel:
		return transfer, nil
	case <-time.After(5 * time.Second):
		return models.Transfer{}, errors.New("persistence operation timed out")
	}
}

func createTransactions(transfer models.Transfer) (models.Transaction, models.Transaction, error) {
	debit := models.Transaction{}
	debit.AccountId = transfer.FromAccount
	debit.Amount = transfer.Amount
	debit.TransactionType = "DEBIT"

	credit := models.Transaction{}
	credit.AccountId = transfer.ToAccount
	credit.Amount = transfer.Amount
	credit.TransactionType = "CREDIT"

	return debit, credit, nil
}
