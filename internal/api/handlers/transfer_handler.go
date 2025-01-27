package handlers

import (
	"encoding/json"
	"goBank/internal/services"
	"net/http"
	"strings"
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

		if newTransfer.ToAccount == newTransfer.FromAccount {
			http.Error(w, "You can't transfer to the same account", http.StatusBadRequest)
			return
		}

		response, err := services.CreateTransfer(newTransfer)
		if err != nil {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
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
