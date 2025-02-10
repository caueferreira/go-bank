package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
)

type KafkaConsumer[T any] struct {
	Consumer *kafka.Consumer
	Topic    string
	Messages chan KafkaEnvelope[T]
}

func NewKafkaConsumer[T any](groupID, topic string) (*KafkaConsumer[T], error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer[T]{
		Consumer: c,
		Topic:    topic,
		Messages: make(chan KafkaEnvelope[T]),
	}, nil
}

func (kc *KafkaConsumer[T]) StartListening() {
	go func() {
		for {
			msg, err := kc.Consumer.ReadMessage(-1)
			if err == nil {
				receivedEnvelope, err := DeserializeKafkaEnvelope[T](msg.Value)
				if err == nil {
					kc.Messages <- *receivedEnvelope
				} else {
					log.Printf("Failed to deserialize message: %v", err)
				}
			} else {
				log.Printf("Consumer error: %v", err)
			}
		}
	}()
}

func (kc *KafkaConsumer[T]) Close() {
	close(kc.Messages)
	kc.Consumer.Close()
}
