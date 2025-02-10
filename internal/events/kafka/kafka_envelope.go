package kafka

import "encoding/json"

type KafkaEnvelope[T any] struct {
	MessageID string `json:"messageId"`
	Message   T      `json:"message"`
}

func (e *KafkaEnvelope[T]) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

func DeserializeKafkaEnvelope[T any](data []byte) (*KafkaEnvelope[T], error) {
	var envelope KafkaEnvelope[T]
	err := json.Unmarshal(data, &envelope)
	if err != nil {
		return nil, err
	}
	return &envelope, nil
}
