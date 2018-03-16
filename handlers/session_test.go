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

func TestSession(t *testing.T) {
	sessMgr := session.NewManager()
	sessHdl := &SessionHandler{Manager: sessMgr}
	srv := httptest.NewServer(sessHdl)
	defer srv.Close()

	client := srv.Client()
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		t.Fatalf("cookiejar.New(): %v", err)
	}
	client.Jar = jar
	url := fmt.Sprintf("%s/session", srv.URL)

	checkNoUser := func(label string) {
		// GET /session
		resp, err := client.Get(url)
		if err != nil {
			t.Fatalf("client.Get(%s): %v", url, err)
		}
		if got, want := resp.StatusCode, http.StatusUnauthorized; got != want {
			t.Errorf("%s: GET /session: got: %d, want: %d", label, got, want)
		}

		// DELETE /session
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		if err != nil {
			t.Fatalf("http.NewRequest: %v", err)
		}
		resp, err = client.Do(req)
		if err != nil {
			t.Fatalf("client.Get(%s): %v", url, err)
		}
		if got, want := resp.StatusCode, http.StatusUnauthorized; got != want {
			t.Errorf("%s: DELETE /session: got: %d, want: %d", label, got, want)
		}
	}

	// No user has logged in yet.
	checkNoUser("not logged in yet")

	// Login with a bad user.
	resp, err := client.Post(url, "application/json", bytes.NewReader([]byte(`{"name":"bad","password":"bad"}`)))
	if err != nil {
		t.Fatalf("client.Get(%s): %v", url, err)
	}
	if got, want := resp.StatusCode, http.StatusUnauthorized; got != want {
		t.Errorf("bad login: GET /session: got: %d, want: %d", got, want)
	}
	checkNoUser("bad user attempted")

	// Login with a good user.
	resp, err = client.Post(url, "application/json", bytes.NewReader([]byte(`{"name":"test","password":"test"}`)))
	if err != nil {
		t.Fatalf("client.Get(%s): %v", url, err)
	}
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("good login: GET /session: got: %d, want: %d", got, want)
	}
	// GET /session with good session.
	resp, err = client.Get(url)
	if err != nil {
		t.Fatalf("client.Get(%s): %v", url, err)
	}
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("after login: GET /session: got: %d, want: %d", got, want)
	}

	// Logout user.
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		t.Fatalf("http.NewRequest: %v", err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("client.Get(%s): %v", url, err)
	}
	if got, want := resp.StatusCode, http.StatusNoContent; got != want {
		t.Errorf("logout: DELETE /session: got: %d, want: %d", got, want)
	}

	// Check session is no longer valid.
	checkNoUser("after logout")
}
