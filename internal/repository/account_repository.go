package repository

import (
	"errors"
	"github.com/gocql/gocql"
	"goBank/internal/db"
	"goBank/internal/models"
)

func connectCassandra() *gocql.Session {
	cluster := gocql.NewCluster("localhost:9042")
	cluster.Keyspace = "go_bank"
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	return session
}

func SaveAccount(account models.Account) (models.Account, error) {
	//session := connectCassandra()
	//defer session.Close()
	//
	//err := session.Query("INSERT INTO accounts (id, name, email, sort_code, account_number, balance) VALUES (?,?,?,?,?,?)",
	//	account.ID, account.Name, account.Email, account.SortCode, account.Number, account.Balance).Exec()
	//if err != nil {
	//	return models.Account{}, err
	//}

	db.AccountsMutex.Lock()
	db.Accounts[account.ID] = account
	db.AccountsMutex.Unlock()
	return account, nil
}

func GetAccountById(id string) (models.Account, error) {
	//session := connectCassandra()
	//defer session.Close()
	//
	//var account models.Account
	//err := session.Query("SELECT * FROM accounts WHERE id = ?", id).Scan(&account.ID)
	//if err != nil {
	//	return models.Account{}, err
	//}
	db.AccountsMutex.Lock()
	account, _ := db.Accounts[id]
	db.AccountsMutex.Unlock()

	return account, nil
}

func CreditAccount(credit models.Transaction) (models.Account, error) {
	db.AccountsMutex.Lock()
	account, exists := db.Accounts[credit.AccountId]
	if !exists {
		return models.Account{}, errors.New("Account not found")
	}

	account.Balance += credit.Amount
	db.Accounts[account.ID] = account
	db.AccountsMutex.Unlock()
	return account, nil
}
func DebitAccount(debit models.Transaction) (models.Account, error) {
	db.AccountsMutex.Lock()
	account, exists := db.Accounts[debit.AccountId]
	if !exists {
		return models.Account{}, errors.New("account not found")
	}
	if account.Balance < debit.Amount {
		return models.Account{}, errors.New("not enough balance")
	}
	account.Balance -= debit.Amount
	db.Accounts[account.ID] = account
	db.AccountsMutex.Unlock()
	return account, nil
}

func GetAllAccounts() []models.Account {
	db.AccountsMutex.Lock()
	var accounts []models.Account
	for _, account := range db.Accounts {
		accounts = append(accounts, account)
	}
	db.AccountsMutex.Unlock()
	return accounts
}
