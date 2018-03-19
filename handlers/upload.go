package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/groggygopher/oyster/register"
	"github.com/groggygopher/oyster/session"
)

// NewUploadHandler returns a new UploadHandler with the given SessionManager
func NewUploadHandler(man *session.Manager) *UploadHandler {
	return &UploadHandler{manager: man}
}

// UploadHandler handles transactions imports with CSV.
type UploadHandler struct {
	manager *session.Manager
}

// ServeHTTP handles importing the uploaded transactions and returning the status.
func (uh *UploadHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	usr := RequestUser(uh.manager, req)
	if usr == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Unsupported method %s", req.Method)))
		return
	}

	trans, err := register.ReadAllTransactions(req.Body)
	if err != nil {
		log.Printf("error: register.ReadAllTransactions: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("There was an error. No data was imported."))
		return
	}
	imported := usr.ImportTransactions(trans)

	resp := &struct {
		Uploaded int `json:"uploaded"`
		Imported int `json:"imported"`
	}{
		Uploaded: len(trans),
		Imported: imported,
	}

	jsonEnc, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "json encoding error", http.StatusInternalServerError)
		log.Printf("error: json.Marshal(%v): %v", resp, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonEnc)
}
