package main

import (
	"fmt"
	"github.com/go-faker/faker/v4"
	"goBank/internal/models"
	"goBank/internal/services"
	"net/http"
)

import "goBank/internal/api/handlers"

func main() {
	for i := 1; i <= 10; i++ {
		account := models.Account{}
		account.Name = faker.FirstName() + " " + faker.LastName()
		account.Balance = 10000
		account.Email = faker.Email()
		_, err := services.CreateAccount(account)
		if err != nil {
			return
		}
	}

	http.HandleFunc("/", handlers.HomeHandler)

	http.HandleFunc("/accounts", handlers.AccountsHandler)
	http.HandleFunc("/account/", handlers.AccountHandler)

	http.HandleFunc("/credit", handlers.CreditHandler)
	http.HandleFunc("/debit", handlers.DebitHandler)

	http.HandleFunc("/transfers", handlers.HandleTransfers)
	http.HandleFunc("/transfer/", handlers.HandleTransfer)

	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}
