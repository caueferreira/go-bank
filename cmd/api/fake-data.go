package main

import (
	"github.com/go-faker/faker/v4"
	"goBank/internal/models"
	"goBank/internal/services"
	"math/rand"
)

func GenerateData() {
	for i := 1; i <= 10; i++ {
		newAccount := models.CreateAccount{}
		newAccount.Name = faker.FirstName() + " " + faker.LastName()
		newAccount.Email = faker.Email()
		_, err := services.CreateAccount(newAccount)
		if err != nil {
			return
		}
	}

	accounts := services.GetAccounts()

	for _, account := range accounts.Accounts {
		services.CreateCredit(models.CreateTransaction{AccountId: account.ID, Amount: 10000, TransactionType: "CREDIT"})
	}

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
