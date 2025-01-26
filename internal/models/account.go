package models

type Account struct {
	ID        string `json:"id"`
	SortCode  string `json:"sortCode"`
	Number    string `json:"number"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Balance   int    `json:"balance"`
	CreatedAt int64  `json:"createdAt"`
}

type Accounts struct {
	Accounts []Account `json:"accounts"`
}
