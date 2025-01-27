package memory

import (
	"errors"
	"goBank/internal/db"
	"goBank/internal/models"
)

func SaveTransfer(transfer models.Transfer) (models.Transfer, error) {
	db.TransfersMutex.Lock()
	db.Transfers[transfer.ID] = transfer
	db.TransfersMutex.Unlock()
	return transfer, nil
}

func GetTransferById(id string) (models.Transfer, error) {
	db.TransfersMutex.Lock()
	transfer, exists := db.Transfers[id]
	db.TransfersMutex.Unlock()

	if !exists {
		return models.Transfer{}, errors.New("transfer does not exist")
	}

	return transfer, nil
}

func GetAllTransfers() []models.Transfer {
	db.TransfersMutex.Lock()
	var transfers []models.Transfer
	for _, transfer := range db.Transfers {
		transfers = append(transfers, transfer)
	}
	db.TransfersMutex.Unlock()
	return transfers
}
