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

		select {
		case response := <-kafka.CreateTransferResponseConsumer.Messages:
			if response.MessageID == messageID {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response.Message)
				return
			}
		case <-time.After(5 * time.Second):
			http.Error(w, "Transfer request timed out", http.StatusGatewayTimeout)
			return
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
