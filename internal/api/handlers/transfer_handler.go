package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"goBank/internal/events/kafka"
	"goBank/internal/services"
	"net/http"
	"strings"
	"time"
)

import "goBank/internal/models"

func HandleTransfers(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var newTransfer models.Transfer
		err := json.NewDecoder(r.Body).Decode(&newTransfer)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		messageID := uuid.New().String()
		envelope := kafka.KafkaEnvelope[models.Transfer]{MessageID: messageID, Message: newTransfer}

		err = kafka.CreateTransferProducer.ProduceMessage(envelope)
		if err != nil {
			http.Error(w, "failed to send transfer request", http.StatusInternalServerError)
			return
		}

		//for {
		//	messages := kafka.CreateTransferResponseConsumer.Messages
		//	for message := range messages {
		//		if message.MessageID == messageID {
		//			w.Header().Set("Content-Type", "application/json")
		//			json.NewEncoder(w).Encode(message.Message)
		//
		//			return
		//		}
		//	}
		//}

		timeout := time.Now().Add(5 * time.Second)
		kafka.TransfersResponseCache.Lock()
		defer kafka.TransfersResponseCache.Unlock()
		for {
			if response, exists := kafka.TransfersResponseCache.PendingResponses[messageID]; exists {
				delete(kafka.TransfersResponseCache.PendingResponses, messageID)

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				return
			}

			if time.Now().After(timeout) {
				http.Error(w, "Transfer request timed out", http.StatusGatewayTimeout)
				return
			}
			kafka.TransfersResponseCache.Cond.Wait()
		}
	} else if r.Method == http.MethodGet {
		response := services.GetTransfers()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/transfer/")
	if id == "" {
		http.Error(w, "Missing account ID", http.StatusBadRequest)
		return
	}

	response, err := services.GetTransferById(id)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
