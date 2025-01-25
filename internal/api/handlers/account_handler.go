package handlers

import (
	"encoding/json"
	"goBank/internal/services"
	"net/http"
	"strings"
)

import "goBank/internal/models"

func AccountsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var newAccount models.Account
		err := json.NewDecoder(r.Body).Decode(&newAccount)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		createdAccount, err := services.CreateAccount(newAccount)
		if err != nil {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(createdAccount)
	} else if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		response := services.GetAccounts()
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func AccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/account/")
	if id == "" {
		http.Error(w, "Missing account ID", http.StatusBadRequest)
		return
	}

	response, err := services.GetAccountById(id)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
