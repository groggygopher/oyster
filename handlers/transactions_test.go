package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"

	"github.com/groggygopher/oyster/session"
	"golang.org/x/net/publicsuffix"
)

func TestGetTransactions(t *testing.T) {
	m, err := session.CreateTestManager()
	if err != nil {
		t.Fatalf("CreateTestManager: %v", err)
	}

	sessHdl := NewSessionHandler(m)
	srv := httptest.NewServer(sessHdl)
	defer srv.Close()

	client := srv.Client()
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		t.Fatalf("cookiejar.New(): %v", err)
	}
	client.Jar = jar
	url := fmt.Sprintf("%s/transactions", srv.URL)

	resp, err := client.Get(url)
	if err != nil {
		t.Fatalf("client.Get(%s): %v", url, err)
	}
	if got, want := resp.StatusCode, http.StatusUnauthorized; got != want {
		t.Fatalf("no login: GET /transactions: got: %d, want: %d", got, want)
	}

	resp, err = client.Post(srv.URL+"/session", "application/json", bytes.NewReader([]byte(`{"name":"test","password":"test"}`)))
	if err != nil {
		t.Fatalf("client.Post(%s): %v", url, err)
	}
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Fatalf("login: POST /session: got: %d, want: %d", got, want)
	}

	resp, err = client.Get(url)
	if err != nil {
		t.Fatalf("client.Get(%s): %v", url, err)
	}
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Fatalf("after login: GET /transactions: got: %d, want: %d", got, want)
	}
}
