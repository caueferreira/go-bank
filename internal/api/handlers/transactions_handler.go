package handlers

import (
	"encoding/json"
	"goBank/internal/services"
	"net/http"
)

func HandleTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := services.GetTransactions()
	json.NewEncoder(w).Encode(response)
}
