package handlers

import (
	"encoding/json"
	"goBank/internal/events"
	"net/http"
	"time"
)

func HandleTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	events.GetAllTransactionsChannel <- struct{}{}

	select {
	case transactions := <-events.GetAllTransactionsResponseChannel:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(transactions)
	case <-time.After(5 * time.Second):
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusGatewayTimeout)

		response := map[string]string{"error": "persistence operation timed out"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		}
	}
}
