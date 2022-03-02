package auth

import (
	"crypto/sha256"
	"fmt"
	"time"
)

type Token struct {
	Value   string    `json:"value"`
	Expires time.Time `json:"expires"`
}

type Handler interface {
	Login(username, password string) (Token, error)
	Authenticate(token string) error
}

func HashSHA256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}
