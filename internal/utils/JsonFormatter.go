package utils

import (
	"encoding/json"
)

func Serialize[T any](message T) ([]byte, error) {
	return json.Marshal(message)
}

func DeserializeMessageData[T any](data []byte) (T, error) {
	var message T
	err := json.Unmarshal(data, message)
	if err != nil {
		return message, err
	}
	return message, nil
}
