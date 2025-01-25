package handlers

import (
	"encoding/json"
	"goBank/internal/services"
	"net/http"
)

import "goBank/internal/models"

func CreditHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var credit models.Transaction
	err := json.NewDecoder(r.Body).Decode(&credit)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	response, err := services.CreateCredit(credit)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
