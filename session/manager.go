package session

import (
	"sync"
	"time"
)

// NewManager returns a new instance of Manager.
func NewManager() *Manager {
	return &Manager{
		active: make(map[string]*Session),
	}
}

// Manager manages the user models and their active state.
type Manager struct {
	sync.Mutex

	active map[string]*Session
}

// ValidSession given a session token produced by Login(), returns the
// associated User with that token value.
func (m *Manager) ValidSession(enc string) *User {
	m.Lock()
	defer m.Unlock()
	if sess, ok := m.active[enc]; ok {
		return sess.User
	}
	return nil
}

// Login given a user name and password, returns a non-nil User pointer and an
// encoded session token when a valid User has logged in. Otherwise, (nil, "")
// is returned.
func (m *Manager) Login(name, password string) (*User, string) {
	if name == "test" && password == "test" {
		usr := &User{
			Name: name,
		}
		s := &Session{
			Start: time.Now(),
			User:  usr,
		}
		enc := s.String()
		m.Lock()
		defer m.Unlock()
		m.active[enc] = s
		return usr, enc
	}
	return nil, ""
}

// Logout removes the given session token and returns true iff the session
// state was modified. A false return value indicates the session token was
// invalid.
func (m *Manager) Logout(enc string) bool {
	m.Lock()
	defer m.Unlock()
	_, ok := m.active[enc]
	delete(m.active, enc)
	return ok
}
