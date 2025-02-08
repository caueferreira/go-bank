package main

import (
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"goBank/internal/events/kafka"
	"goBank/internal/models"
	"goBank/internal/services"
	"math/rand"
	"time"
)

func GenerateData() {
	for i := 1; i <= 10; i++ {
		newAccount := models.CreateAccount{}
		newAccount.Name = faker.FirstName() + " " + faker.LastName()
		newAccount.Email = faker.Email()
		services.CreateAccount(newAccount)
	}

	accounts := services.GetAccounts()

	for _, account := range accounts.Accounts {
		services.CreateCredit(models.CreateTransaction{AccountId: account.ID, Amount: 10000, TransactionType: "CREDIT"})
	}

	for i := 1; i <= 10; i++ {
		rand.NewSource(time.Now().UnixNano())

		fromAccount := accounts.Accounts[rand.Intn(10)]
		toAccount := accounts.Accounts[rand.Intn(10)]

		transfer := models.Transfer{}
		transfer.Amount = rand.Intn(100)
		transfer.ToAccount = toAccount.ID
		transfer.FromAccount = fromAccount.ID

		envelope := kafka.KafkaEnvelope[models.Transfer]{
			MessageID: uuid.New().String(),
			Message:   transfer,
		}

		if fromAccount != toAccount && envelope != envelope {
			//go kafka.CreateTransferProducer.ProduceMessage(envelope)
		}

	}
}
