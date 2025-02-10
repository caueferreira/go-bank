package kafka

import (
	"goBank/internal/models"
	"goBank/internal/services"
	"log"
	"sync"
)

var CreateTransferProducer *KafkaProducer[models.Transfer]
var CreateTransferResponseProducer *KafkaProducer[models.Transfer]

var CreateTransferConsumer *KafkaConsumer[models.Transfer]
var CreateTransferResponseConsumer *KafkaConsumer[models.Transfer]

var TransfersResponseCache = struct {
	sync.Mutex
	Cond             *sync.Cond
	PendingResponses map[string]models.Transfer
}{
	//Cond:             sync.NewCond(&sync.Mutex{}),
	PendingResponses: make(map[string]models.Transfer),
}

func InitTransferKafkaHandlers() {
	TransfersResponseCache.Cond = sync.NewCond(&TransfersResponseCache.Mutex)

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
}

func KafkaTransferCreateWorker(consumer *KafkaConsumer[models.Transfer], producer *KafkaProducer[models.Transfer]) {
	go func() {
		for msg := range consumer.Messages {
			result, err := services.CreateTransfer(msg.Message)
			if err != nil {
				log.Printf("failed to create transfer: %v", err)
			} else {
				producer.ProduceMessage(KafkaEnvelope[models.Transfer]{
					MessageID: msg.MessageID,
					Message:   result,
				})
			}

			TransfersResponseCache.Lock()
			TransfersResponseCache.PendingResponses[msg.MessageID] = result
			TransfersResponseCache.Unlock()
			TransfersResponseCache.Cond.Broadcast()
		}
	}()
}
