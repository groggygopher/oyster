package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/groggygopher/oyster/rule"
	"github.com/groggygopher/oyster/session"
)

// NewRuleHandler returns a new RuleHandler with the given SessionManager
func NewRuleHandler(man *session.Manager) *RuleHandler {
	return &RuleHandler{manager: man}
}

// RuleHandler handles transactions imports with CSV.
type RuleHandler struct {
	manager *session.Manager
}

func (rh *RuleHandler) get(w http.ResponseWriter, req *http.Request) {
	usr := RequestUser(rh.manager, req)
	if usr == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	if err := enc.Encode(usr.RuleManager().Rules()); err != nil {
		log.Printf("error: json.Encode: %v", err)
		return
	}
}

func readRule(r io.Reader) (*rule.Rule, error) {
	rule := &rule.Rule{}
	dec := json.NewDecoder(r)
	if err := dec.Decode(rule); err != nil {
		return nil, fmt.Errorf("json.Decode: %v", err)
	}
	return rule, nil
}

func (rh *RuleHandler) post(w http.ResponseWriter, req *http.Request) {
	usr := RequestUser(rh.manager, req)
	if usr == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	r, err := readRule(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid JSON Rule object"))
		log.Printf("error: readRule: %v", err)
		return
	}
	added := usr.RuleManager().AddRule(r)
	if !added {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("A rule with name '%s' already exists", r.Name)))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (rh *RuleHandler) put(w http.ResponseWriter, req *http.Request) {
	usr := RequestUser(rh.manager, req)
	if usr == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	r, err := readRule(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid JSON Rule object"))
		log.Printf("error: readRule: %v", err)
		return
	}
	usr.RuleManager().UpsertRule(r.Name, r)
	w.WriteHeader(http.StatusNoContent)
}

func (rh *RuleHandler) delete(w http.ResponseWriter, req *http.Request) {
	usr := RequestUser(rh.manager, req)
	if usr == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	r, err := readRule(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid JSON Rule object"))
		log.Printf("error: readRule: %v", err)
		return
	}
	removed := usr.RuleManager().DeleteRule(r.Name)
	if !removed {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("No rule with name '%s' exists", r.Name)))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ServeHTTP handles GET, PUT, POST, and DELETE rule requests.
func (rh *RuleHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	switch req.Method {
	case http.MethodGet:
		rh.get(w, req)
	case http.MethodPost:
		rh.post(w, req)
	case http.MethodPut:
		rh.put(w, req)
	case http.MethodDelete:
		rh.delete(w, req)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("Unsupported method: %s", req.Method)))
	}
}
