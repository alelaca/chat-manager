package jwt

import (
	"net/http"
	"time"

	"github.com/alelaca/chat-manager/src/apperrors"
	"github.com/alelaca/chat-manager/src/auth"
	"github.com/dgrijalva/jwt-go"
)

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

var jwtKey = []byte("my_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Handler struct {
}

func (h *Handler) Login(username, password string) (auth.Token, error) {
	expectedPassword, ok := users[username]
	if !ok || expectedPassword != password {
		return auth.Token{}, apperrors.CreateAPIError(http.StatusUnauthorized, "invlaid username or password")
	}

	tokenExpirationTime := time.Now().Add(15 * time.Minute)

	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return auth.Token{}, apperrors.CreateAPIError(http.StatusInternalServerError, "error processing authentication")
	}

	return auth.Token{
		Value:   tokenString,
		Expires: tokenExpirationTime,
	}, nil
}

func (h *Handler) Authenticate(token string) error {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return apperrors.CreateAPIError(http.StatusUnauthorized, "request anouthorized")
		}

		return apperrors.CreateAPIError(http.StatusBadRequest, "couldn't verify identity")
	}

	if !tkn.Valid {
		return apperrors.CreateAPIError(http.StatusUnauthorized, "request anouthorized")
	}

	return nil
}
