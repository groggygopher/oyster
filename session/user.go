package session

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/groggygopher/oyster/register"
	"github.com/groggygopher/oyster/rule"
)

type serializeableUser struct {
	Name         string
	Transactions []*register.Transaction
	Rules        []*rule.Rule
}

// DeserializeUser takes the given bytes and decodes a User.
func DeserializeUser(b []byte) (*User, error) {
	r := bytes.NewReader(b)
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("gzip.NewReader: %v", err)
	}
	serUsr := &serializeableUser{}
	jsonDec := json.NewDecoder(zr)
	if err := jsonDec.Decode(serUsr); err != nil {
		return nil, fmt.Errorf("json.Decode: %v", err)
	}

	usr := &User{
		Name:         serUsr.Name,
		transactions: serUsr.Transactions,
		manager:      rule.NewManager(serUsr.Rules),
	}
	return usr, nil
}

// User is contains all of a User's data.
type User struct {
	sync.Mutex

	Name string `json:"name"`

	passkey []byte
	// Most recent is at index len() - 1.
	transactions []*register.Transaction
	manager      *rule.Manager
}

// ImportTransactions imports new transactions from the given data, returning
// the number of imported transactions.
func (u *User) ImportTransactions(newTrans []*register.Transaction) int {
	u.Lock()
	defer u.Unlock()

	has := make(map[string]*register.Transaction)
	for _, t := range u.transactions {
		has[t.ID] = t
	}

	var count int
	for _, t := range newTrans {
		if _, ok := has[t.ID]; !ok {
			u.transactions = append(u.transactions, t)
			count++
		}
	}
	sort.SliceStable(u.transactions, func(i, j int) bool {
		return u.transactions[i].Date.After(*u.transactions[j].Date)
	})
	return count
}

// Transactions returns a slice of all transactions for this user.
func (u *User) Transactions() []*register.Transaction {
	u.Lock()
	defer u.Unlock()
	return u.transactions
}

// Serialize generates a binary serialization of this User.
func (u *User) Serialize() ([]byte, error) {
	u.Lock()
	defer u.Unlock()

	serUsr := &serializeableUser{
		Name:         u.Name,
		Transactions: u.transactions,
		Rules:        u.manager.Rules(),
	}
	buf := &bytes.Buffer{}
	zw := gzip.NewWriter(buf)
	jsonEnc := json.NewEncoder(zw)
	if err := jsonEnc.Encode(serUsr); err != nil {
		return nil, fmt.Errorf("json.Encode: %v", err)
	}
	zw.Close()
	return buf.Bytes(), nil
}

// Session is a data model containing the user session state.
type Session struct {
	Start time.Time `json:"start"`
	User  *User     `json:"user"`
}

// String creates a unique session token.
func (s *Session) String() string {
	buf := &bytes.Buffer{}
	zw := gzip.NewWriter(buf)
	jsonEnc := json.NewEncoder(zw)
	if err := jsonEnc.Encode(s); err != nil {
		log.Printf("error: json.Encode(session): %v", err)
		return ""
	}
	zw.Close()
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}
