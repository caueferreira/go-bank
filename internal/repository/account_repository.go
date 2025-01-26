package repository

import (
	"errors"
	"goBank/internal/db"
	"goBank/internal/models"
	"log"
)

func SaveAccount(account models.Account) (models.Account, error) {
	session := db.ConnectCassandra()
	defer session.Close()

	err := session.Query("INSERT INTO accounts (id, name, email, sort_code, account_number, balance, created_at) VALUES (?,?,?,?,?,?, ?)",
		account.ID, account.Name, account.Email, account.SortCode, account.Number, account.Balance, account.CreatedAt).Exec()
	if err != nil {
		log.Fatal(err)
		return models.Account{}, err
	}
	return account, nil
}

func GetAccountById(accountId string) (models.Account, error) {
	session := db.ConnectCassandra()
	defer session.Close()

	var account models.Account

	err := session.Query("SELECT * FROM accounts WHERE id = ?", accountId).Scan(
		&account.ID,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
		&account.Email,
		&account.Name,
		&account.SortCode)

	if err != nil {
		log.Fatal(err)
		return models.Account{}, err
	}

	return account, nil
}

func CreditAccount(credit models.Transaction) (models.Account, error) {
	session := db.ConnectCassandra()
	defer session.Close()

	account, _ := GetAccountById(credit.AccountId)
	err := session.Query("UPDATE accounts SET balance = ? WHERE id = ?", account.Balance+credit.Amount, credit.AccountId).Exec()
	if err != nil {
		log.Fatal(err)
		return models.Account{}, err
	}

	return GetAccountById(credit.AccountId)
}
func DebitAccount(debit models.Transaction) (models.Account, error) {
	session := db.ConnectCassandra()
	defer session.Close()

	account, _ := GetAccountById(debit.AccountId)
	if account.Balance < debit.Amount {
		return models.Account{}, errors.New("not enough balance")
	}

	err := session.Query("UPDATE accounts SET balance = ? WHERE id = ?", account.Balance-debit.Amount, debit.AccountId).Exec()
	if err != nil {
		log.Fatal(err)
		return models.Account{}, err
	}
	return GetAccountById(debit.AccountId)
}

func GetAllAccounts() []models.Account {
	session := db.ConnectCassandra()
	defer session.Close()

	var accounts []models.Account
	iter := session.Query("SELECT * FROM accounts").Iter()

	var id, name, email, sortCode, accountNumber string
	var balance int
	var createdAt int64

	for iter.Scan(&id, &accountNumber, &balance, &createdAt, &email, &sortCode, &name) {
		account := models.Account{
			ID:        id,
			Name:      name,
			Email:     email,
			SortCode:  sortCode,
			Number:    accountNumber,
			Balance:   balance,
			CreatedAt: createdAt,
		}
		accounts = append(accounts, account)
	}

	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}

	return accounts
}
