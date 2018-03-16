package session

import (
	"testing"
	"time"
)

func TestLogout(t *testing.T) {
	m := NewManager()
	m.active["test"] = &Session{
		Start: time.Now(),
		User: &User{
			Name: "test",
		},
	}
	if got := m.Logout("bad"); got {
		t.Error("expected false from unknown session")
	}
	if got := m.Logout("test"); !got {
		t.Error("expected true from known session")
	}
	if got := m.Logout("test"); got {
		t.Error("expected false from logged out session")
	}
	if got := m.ValidSession("test"); got != nil {
		t.Error("expected nil from logged out session")
	}
}

func TestValidSession(t *testing.T) {
	m := NewManager()
	m.active["test"] = &Session{
		Start: time.Now(),
		User: &User{
			Name: "test",
		},
	}
	if got := m.ValidSession("test"); got == nil {
		t.Error("expected non-nil from known session")
	}
	if got := m.ValidSession("bad"); got != nil {
		t.Error("expected nil from unknown session")
	}
}
