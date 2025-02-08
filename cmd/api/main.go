package main

import (
	"fmt"
	"net/http"
)

import "goBank/internal/api/handlers"

func main() {
	InnitDatabase()
	StartWorkers()
	GenerateData()
	http.HandleFunc("/", handlers.HomeHandler)

	http.HandleFunc("/accounts", handlers.AccountsHandler)
	http.HandleFunc("/account/", handlers.AccountHandler)

	http.HandleFunc("/credit", handlers.CreditHandler)
	http.HandleFunc("/debit", handlers.DebitHandler)

	http.HandleFunc("/transfers", handlers.HandleTransfers)
	http.HandleFunc("/transfer/", handlers.HandleTransfer)

	http.HandleFunc("/transactions", handlers.HandleTransactions)

	fmt.Println("Server running on port 8081")
	http.ListenAndServe(":8081", nil)
}
