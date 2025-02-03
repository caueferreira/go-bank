package main

import (
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"goBank/internal/events"
	"goBank/internal/events/kafka"
	"goBank/internal/models"
	"goBank/internal/services"
	"log"
	"math/rand"
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

	producer, err := kafka.NewKafkaProducer[models.Transfer]("transfers")
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	consumer, err := kafka.NewKafkaConsumer[models.Transfer]("transfers", "transfers")
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	go consumer.StartListening()
	events.KafkaTransferCreateWorker(consumer)
	//defer consumer.Close()

	for i := 1; i <= 10000; i++ {
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

		go producer.ProduceMessage(envelope)
		//go services.CreateTransfer(transfer)
	}
}
