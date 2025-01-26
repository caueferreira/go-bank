package cassandra

import (
	"goBank/internal/db"
	"goBank/internal/models"
	"log"
)

func SaveTransaction(transaction models.Transaction) (models.Transaction, error) {
	session := db.ConnectCassandra()
	defer session.Close()

	err := session.Query("INSERT INTO transactions (id, account_id, transaction_tyoe, amount, created_at) VALUES (?,?,?,?,?)",
		transaction.ID, transaction.AccountId, transaction.TransactionType, transaction.Amount, transaction.CreatedAt).Exec()
	if err != nil {
		log.Fatal(err)
		return models.Transaction{}, err
	}
	return transaction, nil
}

func FindTransactionById(transactionId string) (models.Transaction, error) {
	session := db.ConnectCassandra()
	defer session.Close()

	var transaction models.Transaction

	err := session.Query("SELECT * FROM accounts WHERE id = ?", transactionId).Scan(
		&transaction.ID,
		&transaction.AccountId,
		&transaction.TransactionType,
		&transaction.Amount,
		&transaction.CreatedAt)

	if err != nil {
		log.Fatal(err)
		return models.Transaction{}, err
	}

	return transaction, nil
}

func GetAllTransactions() []models.Transaction {
	session := db.ConnectCassandra()
	defer session.Close()

	var transactions []models.Transaction
	iter := session.Query("SELECT * FROM transactions").Iter()

	var id, accountId, transactionType string
	var amount int
	var createdAt int64

	for iter.Scan(&id, &accountId, &amount, &createdAt, &transactionType) {
		transaction := models.Transaction{
			ID:              id,
			AccountId:       accountId,
			Amount:          amount,
			CreatedAt:       createdAt,
			TransactionType: transactionType,
		}
		transactions = append(transactions, transaction)
	}

	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}

	return transactions
}
