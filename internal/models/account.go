package models

type CreateAccount struct {
	RequestId string `json:"requestId"`
	Name      string `json:"name"`
	Email     string `json:"email"`
}

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
