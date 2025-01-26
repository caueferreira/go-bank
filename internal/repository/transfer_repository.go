package repository

import (
	"goBank/internal/db"
	"goBank/internal/models"
	"log"
)

func SaveTransfer(transfer models.Transfer) (models.Transfer, error) {
	session := db.ConnectCassandra()
	defer session.Close()

	err := session.Query("INSERT INTO transfers (id, from_account, to_account, amount, success, created_at) VALUES (?,?,?,?,?,?)",
		transfer.ID, transfer.FromAccount, transfer.ToAccount, transfer.Amount, transfer.Success, transfer.CreatedAt).Exec()
	if err != nil {
		log.Fatal(err)
		return models.Transfer{}, err
	}
	return transfer, nil
}

func GetTransferById(transferId string) (models.Transfer, error) {
	session := db.ConnectCassandra()
	defer session.Close()

	var transfer models.Transfer

	err := session.Query("SELECT * FROM transfers WHERE id = ?", transferId).Scan(
		transfer.ID,
		transfer.FromAccount,
		transfer.ToAccount,
		transfer.Amount,
		transfer.Success,
		transfer.CreatedAt)

	if err != nil {
		log.Fatal(err)
		return models.Transfer{}, err
	}

	return transfer, nil
}

func GetAllTransfers() []models.Transfer {
	session := db.ConnectCassandra()
	defer session.Close()

	var transfers []models.Transfer
	iter := session.Query("SELECT * FROM transfers").Iter()

	var id, fromAccount, toAccount string
	var amount int
	var createdAt int64
	var success bool

	for iter.Scan(&id, &amount, &createdAt, &fromAccount, &success, &toAccount) {
		transfer := models.Transfer{
			ID:          id,
			FromAccount: fromAccount,
			ToAccount:   toAccount,
			Amount:      amount,
			CreatedAt:   createdAt,
			Success:     success,
		}
		transfers = append(transfers, transfer)
	}

	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}

	return transfers
}
