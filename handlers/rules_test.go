package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/groggygopher/oyster/session"
	"golang.org/x/net/publicsuffix"
)

func TestRules(t *testing.T) {
	m, err := session.CreateTestManager()
	if err != nil {
		t.Fatalf("CreateTestManager: %v", err)
	}

	ruleHdl := NewRuleHandler(m)
	srv := httptest.NewServer(ruleHdl)
	defer srv.Close()

	client := srv.Client()
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		t.Fatalf("cookiejar.New(): %v", err)
	}
	client.Jar = jar
	urlStr := fmt.Sprintf("%s/rules", srv.URL)
	u, err := url.Parse(urlStr)
	if err != nil {
		t.Fatalf("url.Parse(%s): %v", urlStr, err)
	}

	// Login.
	_, token, err := m.Login("test", "test")
	if err != nil {
		t.Fatalf("testManager.Login: %v", err)
	}
	jar.SetCookies(u, []*http.Cookie{&http.Cookie{Name: sessCookieKey, Value: token}})

	tests := []struct {
		method   string
		body     string
		wantCode int
	}{
		// Order matters!
		{
			method:   http.MethodGet,
			wantCode: http.StatusOK,
		},
		{
			method:   http.MethodPost,
			body:     `{"}`,
			wantCode: http.StatusBadRequest,
		},
		{
			method:   http.MethodPost,
			body:     `{"Name":"test"}`,
			wantCode: http.StatusNoContent,
		},
		{
			method:   http.MethodPost,
			body:     `{"Name":"test"}`,
			wantCode: http.StatusBadRequest,
		},
		{
			method:   http.MethodPut,
			body:     `{"Name":"test"}`,
			wantCode: http.StatusNoContent,
		},
		{
			method:   http.MethodDelete,
			body:     `{"Name":"test"}`,
			wantCode: http.StatusNoContent,
		},
		{
			method:   http.MethodDelete,
			body:     `{"Name":"test"}`,
			wantCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		req, err := http.NewRequest(test.method, urlStr, bytes.NewReader([]byte(test.body)))
		if err != nil {
			t.Fatalf("http.NewRequest: %v", err)
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("client.Get(%s): %v", urlStr, err)
		}
		if got, want := resp.StatusCode, test.wantCode; got != want {
			t.Errorf("response code: got: %d, want: %d", got, want)
		}

	}
}
