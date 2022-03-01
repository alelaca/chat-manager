package auth

import "time"

type Token struct {
	Value   string    `json:"value"`
	Expires time.Time `json:"expires"`
}

type Handler interface {
	Login(username, password string) (Token, error)
	Authenticate(token string) error
}
