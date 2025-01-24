package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type HelloWorld struct {
	Message string `json:"message"`
}

type Account struct {
	ID       string `json:"id"`
	SortCode string `json:"sortCode"`
	Number   string `json:"number"`
	Name     string `json:"name"`
	Amount   int    `json:"amount"`
}

type Accounts struct {
	Accounts []Account `json:"accounts"`
}

var (
	accounts      = make(map[string]Account)
	accountsMutex = &sync.Mutex{}
)

type Transfer struct {
	ID          string `json:"id"`
	ToAccount   string `json:"toAccount"`
	FromAccount string `json:"fromAccount"`
	Amount      int    `json:"amount"`
	Success     bool   `json:"success"`
}

type Transfers struct {
	Transfers []Transfer `json:"transfers"`
}

var (
	transfers      = make(map[string]Transfer)
	transfersMutex = &sync.Mutex{}
)

type Credit struct {
	ID        string `json:"id"`
	AccountId string `json:"accountId"`
	Amount    int    `json:"amount"`
}

type Debit struct {
	ID        string `json:"id"`
	AccountId string `json:"accountId"`
	Amount    int    `json:"amount"`
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	helloWorld := &HelloWorld{}
	helloWorld.Message = "Hello GoBank!"

	json.NewEncoder(w).Encode(helloWorld)
}

func accountsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var newAccount Account
		err := json.NewDecoder(r.Body).Decode(&newAccount)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		newAccount.ID = uuid.New().String()
		newAccount.SortCode = "001942"
		newAccount.Number = strconv.Itoa(10000000 + rand.Intn(99999999-10000000))
		accountsMutex.Lock()
		accounts[newAccount.ID] = newAccount
		accountsMutex.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newAccount)
	} else if r.Method == http.MethodGet {
		accountsMutex.Lock()
		var accountList []Account
		for _, account := range accounts {
			accountList = append(accountList, account)
		}
		accountsResponse := Accounts{}
		accountsResponse.Accounts = accountList
		accountsMutex.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(accountsResponse)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func accountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/account/")
	if id == "" {
		http.Error(w, "Missing account ID", http.StatusBadRequest)
		return
	}

	accountsMutex.Lock()
	account, exists := accounts[id]
	accountsMutex.Unlock()

	if !exists {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

func creditHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var credit Credit
	err := json.NewDecoder(r.Body).Decode(&credit)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	credit.ID = uuid.New().String()

	accountsMutex.Lock()
	account, exists := accounts[credit.AccountId]
	account.Amount += credit.Amount
	accounts[account.ID] = account
	accountsMutex.Unlock()

	if !exists {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

func debitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var debit Debit
	err := json.NewDecoder(r.Body).Decode(&debit)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	debit.ID = uuid.New().String()

	accountsMutex.Lock()
	account, exists := accounts[debit.AccountId]
	account.Amount -= debit.Amount
	accounts[account.ID] = account
	accountsMutex.Unlock()

	if !exists {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

func handleTransfers(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		var newTransfer Transfer
		err := json.NewDecoder(r.Body).Decode(&newTransfer)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		newTransfer.ID = uuid.New().String()
		newTransfer.Success = false

		accountsMutex.Lock()
		toAccount, toExists := accounts[newTransfer.ToAccount]
		fromAccount, fromExists := accounts[newTransfer.FromAccount]

		if !toExists || !fromExists {
			accountsMutex.Unlock()
			http.Error(w, "Account not found", http.StatusNotFound)
			return
		}

		if fromAccount.Amount >= newTransfer.Amount {
			toAccount.Amount += newTransfer.Amount
			fromAccount.Amount -= newTransfer.Amount
			accounts[toAccount.ID] = toAccount
			accounts[fromAccount.ID] = fromAccount
			newTransfer.Success = true
		}
		accountsResponse := Accounts{Accounts: []Account{toAccount, fromAccount}}
		accountsMutex.Unlock()

		transfersMutex.Lock()
		transfers[newTransfer.ID] = newTransfer
		transfersMutex.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(accountsResponse)
	} else if r.Method == http.MethodGet {
		transfersMutex.Lock()
		var transferList []Transfer
		for _, transfer := range transfers {
			transferList = append(transferList, transfer)
		}
		transfersResponse := Transfers{}
		transfersResponse.Transfers = transferList
		transfersMutex.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(transfersResponse)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/transfer/")
	if id == "" {
		http.Error(w, "Missing account ID", http.StatusBadRequest)
		return
	}

	transfersMutex.Lock()
	transfer, exists := transfers[id]
	transfersMutex.Unlock()

	if !exists {
		http.Error(w, "Transfer not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transfer)
}

func main() {
	http.HandleFunc("/", homeHandler)

	http.HandleFunc("/accounts", accountsHandler)
	http.HandleFunc("/account/", accountHandler)

	http.HandleFunc("/credit", creditHandler)
	http.HandleFunc("/debit", debitHandler)

	http.HandleFunc("/transfers", handleTransfers)
	http.HandleFunc("/transfer/", handleTransfer)

	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}
