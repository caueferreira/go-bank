package handlers

import (
	"encoding/json"
	"goBank/internal/services"
	"net/http"
)

import "goBank/internal/models"

func DebitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var debit models.Debit
	err := json.NewDecoder(r.Body).Decode(&debit)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	response, err := services.CreateDebit(debit)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
