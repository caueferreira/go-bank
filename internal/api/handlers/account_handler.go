package handlers

import (
	"encoding/json"
	"goBank/internal/events"
	"net/http"
	"strings"
	"time"
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

		events.AccountCreateChannel <- newAccount

		select {
		case account := <-events.AccountResponseChannel:
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
	} else if r.Method == http.MethodGet {
		events.GetAllAccountsChannel <- struct{}{}

		select {
		case accounts := <-events.GetAllAccountsResponseChannel:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(accounts)
		case <-time.After(5 * time.Second):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusGatewayTimeout)

			response := map[string]string{"error": "persistence operation timed out"}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
			}
		}
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

	events.FindAccountChannel <- id
	select {
	case account := <-events.AccountResponseChannel:
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
