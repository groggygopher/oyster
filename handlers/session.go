package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/groggygopher/oyster/session"
)

const sessCookieKey = "session"

func getCookie(req *http.Request) (string, bool) {
	for _, c := range req.Cookies() {
		if c.Name == sessCookieKey {
			return c.Value, true
		}
	}
	return "", false
}

// SessionHandler manages the sessions of all users and provides an HTTP API to
// support those actions.
type SessionHandler struct {
	Manager *session.Manager
}

func (sh *SessionHandler) checkActiveSession(w http.ResponseWriter, req *http.Request) {
	c, ok := getCookie(req)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	usr := sh.Manager.ValidSession(c)
	if usr == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	if err := enc.Encode(usr); err != nil {
		log.Printf("error: encode user %v: %v", usr, err)
	}
}

func (sh *SessionHandler) login(w http.ResponseWriter, req *http.Request) {
	body := &struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}{}
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(body); err != nil {
		http.Error(w, "invalid user request body", http.StatusBadRequest)
		log.Printf("error: decode user request body: %v", err)
		return
	}
	usr, c := sh.Manager.Login(body.Name, body.Password)
	if usr == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  sessCookieKey,
		Value: c,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	if err := enc.Encode(usr); err != nil {
		log.Printf("error: encode user %v: %v", usr, err)
	}
}

func (sh *SessionHandler) logout(w http.ResponseWriter, req *http.Request) {
	c, ok := getCookie(req)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if ok := sh.Manager.Logout(c); !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: sessCookieKey})
	w.WriteHeader(http.StatusNoContent)
}

// ServeHTTP provides the HTTP API for session management.
// GET requests validate that a user is still logged in.
// POST requests login a user.
// DELETE requests logout a user.
func (sh *SessionHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	switch req.Method {
	case http.MethodGet:
		sh.checkActiveSession(w, req)
	case http.MethodPost:
		sh.login(w, req)
	case http.MethodDelete:
		sh.logout(w, req)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
