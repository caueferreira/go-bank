package handlers

import (
	"encoding/json"
	"net/http"
)

type HelloWorld struct {
	Message string `json:"message"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	helloWorld := &HelloWorld{}
	helloWorld.Message = "Hello GoBank!"

	json.NewEncoder(w).Encode(helloWorld)
}
