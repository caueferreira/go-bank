package handlers

import (
	"encoding/json"
	"goBank/internal/events"
	"net/http"
	"time"
)

import "goBank/internal/models"

func DebitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newDebit models.CreateTransaction
	err := json.NewDecoder(r.Body).Decode(&newDebit)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	events.TransactionCreateRoutine <- newDebit

	select {
	case account := <-events.TransactionResponseRoutine:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(account)
	case <-time.After(5 * time.Second):
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusGatewayTimeout)

		response := map[string]string{"error": "persistence operation timed out"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		}
	}
}
