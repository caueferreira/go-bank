package cassandra

import (
	"goBank/internal/db"
	"goBank/internal/models"
	"log"
)

func SaveTransaction(transaction models.Transaction) (models.Transaction, error) {
	//db.ConnectCassandra()

	err := db.GetSession().Query("INSERT INTO transactions (id, account_id, transaction_type, amount, created_at) VALUES (?,?,?,?,?)",
		transaction.ID, transaction.AccountId, transaction.TransactionType, transaction.Amount, transaction.CreatedAt).Exec()
	if err != nil {
		log.Fatal("SaveTransaction:" + err.Error())
		return models.Transaction{}, err
	}
	return transaction, nil
}

func FindTransactionById(transactionId string) (models.Transaction, error) {
	//db.ConnectCassandra()

	var transaction models.Transaction

	err := db.GetSession().Query("SELECT * FROM accounts WHERE id = ?", transactionId).Scan(
		&transaction.ID,
		&transaction.AccountId,
		&transaction.TransactionType,
		&transaction.Amount,
		&transaction.CreatedAt)

	if err != nil {
		log.Fatal("FindTransactionById:" + err.Error())
		return models.Transaction{}, err
	}

	return transaction, nil
}

func GetAllTransactions() []models.Transaction {
	//db.ConnectCassandra()

	var transactions []models.Transaction
	iter := db.GetSession().Query("SELECT * FROM transactions").Iter()

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
		log.Fatal("GetAllTransactions:" + err.Error())
	}

	return transactions
}
