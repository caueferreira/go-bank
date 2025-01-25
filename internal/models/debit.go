package models

type Debit struct {
	ID        string `json:"id"`
	AccountId string `json:"accountId"`
	Amount    int    `json:"amount"`
}
