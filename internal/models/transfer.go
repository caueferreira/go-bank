package models

type Transfer struct {
	ID          string `json:"id"`
	ToAccount   string `json:"toAccount"`
	FromAccount string `json:"fromAccount"`
	Amount      int    `json:"amount"`
	Success     bool   `json:"success"`
}

type Transfers struct {
	Transfers []Transfer `json:"transfers"`
}
