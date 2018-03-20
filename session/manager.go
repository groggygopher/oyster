package session

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	postfix   = "-Oyster"
	cost      = 13
	nonceLen  = 12
	aesKeyLen = 32
)

func decodeUser(saveFile string, passkey []byte) (*User, error) {
	file, err := ioutil.ReadFile(saveFile)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadFile(%s): %v", saveFile, err)
	}
	nonce := file[:nonceLen]
	encrypted := file[nonceLen:]

	block, err := aes.NewCipher(passkey)
	if err != nil {
		return nil, fmt.Errorf("aes.NewCipher: %v", err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cipher.NewGCM: %v", err)
	}
	plain, err := aesgcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return nil, fmt.Errorf("aesgcm.Open: %v", err)
	}

	usr, err := DeserializeUser(plain)
	if err != nil {
		return nil, fmt.Errorf("DeserializeUser: %v", err)
	}
	usr.passkey = passkey
	return usr, nil
}

func encodeUser(usr *User, saveFile string) error {
	passkey := usr.passkey
	// Don't serialize the passkey, but add it back before returning.
	usr.passkey = nil
	defer func() {
		usr.passkey = passkey
	}()

	serUsr, err := usr.Serialize()
	if err != nil {
		return fmt.Errorf("user.Serialize: %v", err)
	}

	block, err := aes.NewCipher(passkey)
	if err != nil {
		return fmt.Errorf("aes.NewCipher: %v", err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("cipher.NewGCM: %v", err)
	}
	nonce := make([]byte, nonceLen)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("io.ReadFull(random): %v", err)
	}

	outFile, err := os.Create(saveFile)
	if err != nil {
		return fmt.Errorf("os.Create(%s): %v", saveFile, err)
	}
	defer outFile.Close()

	encrypted := aesgcm.Seal(nil, nonce, serUsr, nil)

	if _, err := outFile.Write(nonce); err != nil {
		return fmt.Errorf("outFile.Write: %v", err)
	}
	if _, err := outFile.Write(encrypted); err != nil {
		return fmt.Errorf("outFile.Write: %v", err)
	}
	return nil
}

func generatePasskey(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	passkey := make([]byte, aesKeyLen)
	for i := 0; i < len(hash); i++ {
		passkey[i%aesKeyLen] = passkey[i%aesKeyLen] ^ hash[i]
	}
	return passkey
}

// NewManager returns a new instance of Manager.
func NewManager(saveDir string) *Manager {
	return &Manager{
		active:  make(map[string]*Session),
		saveDir: saveDir,
	}
}

// Manager manages the user models and their active state.
type Manager struct {
	sync.Mutex

	saveDir string
	active  map[string]*Session
}

func (m *Manager) userSaveFile(name string) string {
	baseName := name + postfix
	baseName = base64.URLEncoding.EncodeToString([]byte(baseName))
	return filepath.Join(m.saveDir, baseName)
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

// Register creates a new User, so long as the given name is not being used by
// any other User account. If it is, (nil, "", err) is returned, otherwise the
// new User struct and a session token is returned with nil error.
func (m *Manager) Register(name, password string) (*User, string, error) {
	if len(name) == 0 {
		return nil, "", errors.New("Username must be at least 1 character long")
	}
	if len(password) < 4 {
		return nil, "", errors.New("Password must be at least 4 characters long")
	}

	m.Lock()
	defer m.Unlock()

	errAlreadyExists := fmt.Errorf("there is already a user with name: %s", name)
	saveFile := m.userSaveFile(name)
	if _, err := os.Stat(saveFile); !os.IsNotExist(err) {
		return nil, "", errAlreadyExists
	}
	for _, sess := range m.active {
		if activeName := sess.User.Name; name == activeName {
			return nil, "", errAlreadyExists
		}
	}

	passkey := generatePasskey(password)
	usr := &User{
		Name:    name,
		passkey: passkey,
	}

	s := &Session{
		Start: time.Now(),
		User:  usr,
	}
	enc := s.String()
	m.active[enc] = s
	return usr, enc, nil
}

// Login given a user name and password, returns a non-nil User pointer and an
// encoded session token when a valid User has logged in. Otherwise, (nil, "")
// is returned.
func (m *Manager) Login(name, password string) (*User, string, error) {
	saveFile := m.userSaveFile(name)
	if _, err := os.Stat(saveFile); os.IsNotExist(err) {
		return nil, "", nil
	}

	passkey := generatePasskey(password)

	usr, err := decodeUser(saveFile, passkey)
	if err != nil {
		return nil, "", fmt.Errorf("decodeUser(%s): %v", saveFile, err)
	}

	s := &Session{
		Start: time.Now(),
		User:  usr,
	}
	enc := s.String()
	m.Lock()
	defer m.Unlock()
	m.active[enc] = s
	return usr, enc, nil
}

// Logout removes the given session token and serializes the user to disk.
func (m *Manager) Logout(enc string) error {
	m.Lock()
	defer m.Unlock()
	sess, ok := m.active[enc]
	if !ok {
		return fmt.Errorf("no active session with: %s", enc)
	}
	defer delete(m.active, enc)

	usr := sess.User
	saveFile := m.userSaveFile(usr.Name)
	if err := encodeUser(usr, saveFile); err != nil {
		return fmt.Errorf("encodeUser: %v", err)
	}
	return nil
}

// Close removes all session state from memory, saving it to disk.
// The Manager should not be used after calling Close().
func (m *Manager) Close() error {
	m.Lock()
	defer m.Unlock()
	var errs []string
	for _, sess := range m.active {
		err := func() error {
			usr := sess.User

			saveFile := m.userSaveFile(usr.Name)
			if err := encodeUser(usr, saveFile); err != nil {
				return fmt.Errorf("encodeUser: %v", err)
			}
			return nil
		}()
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	m.active = nil
	if l := len(errs); l > 0 {
		return fmt.Errorf("%d errors closing manager: %s", l, strings.Join(errs, ", "))
	}
	return nil
}

// CreateTestManager returns a Manager with a single valid user login.
// username: test, password: test
func CreateTestManager() (*Manager, error) {
	saveDir := filepath.Join(os.TempDir(), "oyster-test")
	if err := os.RemoveAll(saveDir); err != nil {
		return nil, fmt.Errorf("os.RemoveAll(%s): %v", saveDir, err)
	}
	if err := os.Mkdir(saveDir, os.FileMode(0775)); err != nil {
		return nil, fmt.Errorf("os.Mkdir(%s): %v", saveDir, err)
	}

	passkey := generatePasskey("test")

	manager := NewManager(saveDir)
	testUsr := &User{
		Name:    "test",
		passkey: passkey,
	}
	if err := encodeUser(testUsr, manager.userSaveFile("test")); err != nil {
		return nil, fmt.Errorf("encodeUser: %v", err)
	}

	return manager, nil
}
