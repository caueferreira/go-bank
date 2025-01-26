package models

import "fmt"

type TransactionType string

const (
	CREDIT TransactionType = "CREDIT"
	DEBIT  TransactionType = "DEBIT"
)

func ValidateTransactionType(tt TransactionType) error {
	switch tt {
	case CREDIT, DEBIT:
		return nil
	default:
		return fmt.Errorf("invalid TransactionType: %s", tt)
	}
}

type Transaction struct {
	ID              string `json:"id"`
	AccountId       string `json:"accountId"`
	Amount          int    `json:"amount"`
	TransactionType string `json:"type"`
	CreatedAt       int64  `json:"createdAt"`
}

type Transactions struct {
	Transactions []Transaction `json:"transactions"`
}
