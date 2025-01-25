package models

type Credit struct {
	ID        string `json:"id"`
	AccountId string `json:"accountId"`
	Amount    int    `json:"amount"`
}
