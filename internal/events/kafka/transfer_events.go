package kafka

import (
	"goBank/internal/models"
	"goBank/internal/services"
	"log"
)

var CreateTransferProducer *KafkaProducer[models.Transfer]
var CreateTransferResponseProducer *KafkaProducer[models.Transfer]

var CreateTransferConsumer *KafkaConsumer[models.Transfer]
var CreateTransferResponseConsumer *KafkaConsumer[models.Transfer]

func InitTransferKafkaHandlers() {
	var err error
	CreateTransferProducer, err = NewKafkaProducer[models.Transfer]("create-transfer")
	if err != nil {
		log.Fatalf("Failed to initialize transfer producer: %v", err)
	}

	CreateTransferResponseProducer, err = NewKafkaProducer[models.Transfer]("create-transfer-response")
	if err != nil {
		log.Fatalf("Failed to initialize transfer producer: %v", err)
	}

	CreateTransferResponseConsumer, err = NewKafkaConsumer[models.Transfer]("create-transfer-response-group", "create-transfer-response")
	if err != nil {
		log.Fatalf("Failed to initialize transfer response consumer: %v", err)
	}
	go CreateTransferResponseConsumer.StartListening()

	CreateTransferConsumer, err = NewKafkaConsumer[models.Transfer]("create-transfer-group", "create-transfer")
	if err != nil {
		log.Fatalf("Failed to initialize transfer response producer: %v", err)
	}
	go CreateTransferConsumer.StartListening()

	go KafkaTransferCreateWorker(CreateTransferConsumer, CreateTransferResponseProducer)
	//go TransferRequestWorker(CreateTransferResponseConsumer, CreateTransferProducer)
}

//func TransferRequestWorker(consumer *KafkaConsumer[models.Transfer], producer *KafkaProducer[models.Transfer]) {
//	go func() {
//		for msg := range consumer.Messages {
//			go func(msg KafkaEnvelope[models.Transfer]) {
//				createdTransfer, err := services.CreateTransfer(msg.Message)
//				if err != nil {
//					log.Printf("error creating transfer: %v", err)
//				}
//
//				responseEnvelope := KafkaEnvelope[models.Transfer]{
//					MessageID: msg.MessageID,
//					Message:   createdTransfer,
//				}
//
//				err = producer.ProduceMessage(responseEnvelope)
//				if err != nil {
//					log.Printf("failed to send transfer response: %v", err)
//				}
//			}(msg)
//		}
//	}()
//}

func KafkaTransferCreateWorker(consumer *KafkaConsumer[models.Transfer], producer *KafkaProducer[models.Transfer]) {
	go func() {
		for msg := range consumer.Messages {
			//go func(msg KafkaEnvelope[models.Transfer]) {
			result, err := services.CreateTransfer(msg.Message)
			if err != nil {
				log.Printf("failed to create transfer: %v", err)
			} else {
				producer.ProduceMessage(KafkaEnvelope[models.Transfer]{
					MessageID: msg.MessageID,
					Message:   result,
				})
			}
			//}(msg)
		}
	}()
}
