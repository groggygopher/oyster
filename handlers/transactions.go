package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/groggygopher/oyster/session"
)

// NewTransactionsHandler returns a new TransactionsHandlers with the given SessionManager.
func NewTransactionsHandler(man *session.Manager) *TransactionsHandler {
	return &TransactionsHandler{manager: man}
}

// TransactionsHandler serves GET queries for returning all of a user's transactions.
type TransactionsHandler struct {
	manager *session.Manager
}

// ServeHTTP serves GET queries for returning all of a user's transactions.
func (th *TransactionsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	usr := RequestUser(th.manager, req)
	if usr == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("Unsupported method %s", req.Method)))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	if err := enc.Encode(usr.Transactions()); err != nil {
		log.Printf("error: json.Marshal: %v", err)
		return
	}
}
