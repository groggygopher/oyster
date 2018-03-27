package session

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/groggygopher/oyster/register"
	"github.com/groggygopher/oyster/rule"
)

func TestGeneratePasskey(t *testing.T) {
	// Passkey generation should be deterministic for the same input.
	words := []string{"hello", "", "PASSWORD"}
	for _, w := range words {
		pk := generatePasskey(w)
		for i := 0; i < 10; i++ {
			if got, want := generatePasskey(w), pk; !reflect.DeepEqual(got, want) {
				t.Errorf("%s: got: %v, want: %v", w, got, want)
			}
		}
	}
}

func TestRegister(t *testing.T) {
	tests := []struct {
		label      string
		wantErr    bool
		saveUser   string
		activeUser string
		newUser    string
		password   string
	}{
		// Watch out: Don't use name=test since the test manager already has that.
		{
			label:      "no conflict",
			saveUser:   "test1",
			activeUser: "test2",
			newUser:    "test3",
			password:   "test",
		},
		{
			label:      "save conflict",
			wantErr:    true,
			saveUser:   "test3",
			activeUser: "test2",
			newUser:    "test3",
			password:   "test",
		},
		{
			label:      "active conflict",
			wantErr:    true,
			saveUser:   "test1",
			activeUser: "test3",
			newUser:    "test3",
			password:   "test",
		},
		{
			label:      "both conflict",
			wantErr:    true,
			saveUser:   "test3",
			activeUser: "test3",
			newUser:    "test3",
			password:   "test",
		},
		{
			label:    "name too short",
			wantErr:  true,
			password: "test",
		},
		{
			label:    "password too short",
			wantErr:  true,
			newUser:  "test3",
			password: "tes",
		},
	}

	for _, test := range tests {
		t.Run(test.label, func(t *testing.T) {
			m, err := CreateTestManager()
			if err != nil {
				t.Fatalf("CreateTestManager: %v", err)
			}

			if len(test.saveUser) > 0 {
				saveFile := m.userSaveFile(test.saveUser)
				if _, err := os.Create(saveFile); err != nil {
					t.Fatalf("os.Create(%s): %v", saveFile, err)
				}
			}

			if len(test.activeUser) > 0 {
				m.active["test"] = &Session{
					Start: time.Now(),
					User: &User{
						Name: test.activeUser,
					},
				}
			}

			_, _, err = m.Register(test.newUser, test.password)
			if got, want := err != nil, test.wantErr; got != want {
				t.Errorf("wantErr: err: %v, got: %t, want: %t", err, got, want)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	m, err := CreateTestManager()
	if err != nil {
		t.Fatalf("CreateTestManager: %v", err)
	}

	m.active["test"] = &Session{
		Start: time.Now(),
		User: &User{
			Name:    "test",
			passkey: []byte("testtesttesttesttesttesttesttest"),
		},
	}
	if err := m.Logout("bad"); err == nil {
		t.Error("expected non-nil error from unknown session")
	}
	if err := m.Logout("test"); err != nil {
		t.Errorf("expected nil error from known session, got: %v", err)
	}
	if err := m.Logout("test"); err == nil {
		t.Error("expected non-nil error from logged out session")
	}
	if got := m.ValidSession("test"); got != nil {
		t.Error("expected nil from logged out session")
	}
}

func TestValidSession(t *testing.T) {
	m, err := CreateTestManager()
	if err != nil {
		t.Fatalf("CreateTestManager: %v", err)
	}

	m.active["test"] = &Session{
		Start: time.Now(),
		User: &User{
			Name:    "test",
			passkey: []byte("testtesttesttesttesttesttesttest"),
		},
	}
	if got := m.ValidSession("test"); got == nil {
		t.Error("expected non-nil from known session")
	}
	if got := m.ValidSession("bad"); got != nil {
		t.Error("expected nil from unknown session")
	}
}

func TestDecodeEncode(t *testing.T) {
	passkey := []byte("testtesttesttesttesttesttesttest")
	usr := &User{
		Name:    "test",
		passkey: passkey,
		transactions: []*register.Transaction{
			&register.Transaction{
				Description: "test",
			},
		},
		manager: rule.NewEmptyManager(),
	}
	saveDir := filepath.Join(os.TempDir(), "oyster-test")
	if err := os.RemoveAll(saveDir); err != nil {
		t.Fatalf("os.RemoveAll(%s): %v", saveDir, err)
	}
	if err := os.Mkdir(saveDir, os.FileMode(0775)); err != nil {
		t.Fatalf("os.Mkdir(%s): %v", saveDir, err)
	}
	saveFile := filepath.Join(saveDir, "test")

	if err := encodeUser(usr, saveFile); err != nil {
		t.Fatalf("encodeUser: %v", err)
	}
	decUsr, err := decodeUser(saveFile, passkey)
	if err != nil {
		t.Fatalf("decodeUser: %v", err)
	}

	if got, want := decUsr, usr; !reflect.DeepEqual(got, want) {
		t.Errorf("got: %v, want: %v", got, want)
	}
}
