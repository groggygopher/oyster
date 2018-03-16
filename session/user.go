package session

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"log"
	"time"
)

// User is a data model suitable for passing to the client.
type User struct {
	Name string `json:"name"`
}

// Session is a data model containing the user session state.
type Session struct {
	Start time.Time `json:"start"`
	User  *User     `json:"user"`
}

// String creates a unique session token.
func (s *Session) String() string {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(s); err != nil {
		log.Printf("error: session.String(): %v", err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}
