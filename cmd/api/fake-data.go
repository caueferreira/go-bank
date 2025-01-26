package main

import (
	"github.com/go-faker/faker/v4"
	"goBank/internal/models"
	"goBank/internal/services"
	"math/rand"
)

func GenerateData() {
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

	accounts := services.GetAccounts()

	for i := 1; i <= 1000; i++ {
		fromAccount := accounts.Accounts[rand.Intn(10)]
		toAccount := accounts.Accounts[rand.Intn(10)]

		transfer := models.Transfer{}
		transfer.Amount = rand.Intn(100)
		transfer.ToAccount = toAccount.ID
		transfer.FromAccount = fromAccount.ID

		services.CreateTransfer(transfer)

	}
}
