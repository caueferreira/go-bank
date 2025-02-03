package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
)

type KafkaProducer[T any] struct {
	Producer *kafka.Producer
	Topic    string
}

func NewKafkaProducer[T any](topic string) (*KafkaProducer[T], error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		return nil, err
	}
	return &KafkaProducer[T]{Producer: p, Topic: topic}, nil
}

func (kp *KafkaProducer[T]) ProduceMessage(envelope KafkaEnvelope[T]) error {
	data, err := envelope.Serialize()
	if err != nil {
		return err
	}

	deliveryChan := make(chan kafka.Event)

	err = kp.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kp.Topic, Partition: kafka.PartitionAny},
		Key:            []byte(envelope.MessageID),
		Value:          data,
	}, deliveryChan)

	if err != nil {
		close(deliveryChan)
		return err
	}

	e := <-deliveryChan
	m, ok := e.(*kafka.Message)
	close(deliveryChan)

	if !ok {
		return err
	}

	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}

	log.Printf("Produced message: MessageID=%s, Message=%+v", envelope.MessageID, envelope.Message)
	return nil
}

func (kp *KafkaProducer[T]) Close() {
	kp.Producer.Flush(15000)
	kp.Producer.Close()
}
