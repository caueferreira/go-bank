package services

import (
	"github.com/google/uuid"
	"goBank/internal/models"
	"goBank/internal/repository/cassandra"
	"time"
)

func CreateTransfer(transfer models.Transfer) (models.Transfer, error) {
	transfer.ID = uuid.New().String()
	transfer.Success = false
	transfer.CreatedAt = time.Now().Unix()

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

	_, creditErr := CreateCredit(newCredit)

	if creditErr != nil {
		refund := models.CreateTransaction{}
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
	return cassandra.GetTransferById(id)
}

func GetTransfers() models.Transfers {
	return models.Transfers{Transfers: cassandra.GetAllTransfers()}
}

func saveTransfer(transfer models.Transfer) (models.Transfer, error) {
	savedTransfer, err := cassandra.SaveTransfer(transfer)
	if err != nil {
		return models.Transfer{}, err
	}

	return savedTransfer, nil
}

func createTransactions(transfer models.Transfer) (models.CreateTransaction, models.CreateTransaction, error) {
	debit := models.CreateTransaction{}
	debit.AccountId = transfer.FromAccount
	debit.Amount = transfer.Amount
	debit.TransactionType = "DEBIT"

	credit := models.CreateTransaction{}
	credit.AccountId = transfer.ToAccount
	credit.Amount = transfer.Amount
	credit.TransactionType = "CREDIT"

	return debit, credit, nil
}
