package services

import (
	"github.com/google/uuid"
	"goBank/internal/models"
	"goBank/internal/repository"
)

func CreateTransfer(transfer models.Transfer) (models.Transfer, error) {
	transfer.ID = uuid.New().String()
	transfer.Success = false

	newDebit, newCredit, err := createTransactions(transfer)
	if err != nil {
		saveTransfer(transfer)
		return models.Transfer{}, err
	}

	_, debitErr := CreateDebit(newDebit)
	if debitErr != nil {
		saveTransfer(transfer)
		return models.Transfer{}, debitErr
	}

	repository.SaveTransaction(newDebit)

	_, creditErr := CreateCredit(newCredit)

	if creditErr != nil {
		refund := models.Transaction{}
		refund.AccountId = transfer.FromAccount
		refund.Amount = transfer.Amount
		refund.TransactionType = "CREDIT"

		_, err := CreateCredit(refund)
		if err != nil {
			saveTransfer(transfer)
			return models.Transfer{}, err
		}
	}

	if debitErr == nil && creditErr == nil {
		transfer.Success = true
	}

	savedTransfer, err := saveTransfer(transfer)
	if err != nil {
		return models.Transfer{}, err
	}

	return savedTransfer, nil
}

func GetTransferById(id string) (models.Transfer, error) {
	return repository.GetTransferById(id)
}

func GetTransfers() models.Transfers {
	return models.Transfers{Transfers: repository.GetAllTransfers()}
}

func saveTransfer(transfer models.Transfer) (models.Transfer, error) {
	savedTransfer, err := repository.SaveTransfer(transfer)
	if err != nil {
		return models.Transfer{}, err
	}

	return savedTransfer, nil
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
