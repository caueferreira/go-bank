package models

type Transaction struct {
	ID        string `json:"id"`
	AccountId string `json:"accountId"`
	Amount    int    `json:"amount"`
	Type      string `json:"type"`
}

type Transactions struct {
	Transactions []Transaction `json:"transactions"`
}
