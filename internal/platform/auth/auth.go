package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/axwilliams/books-api/internal/platform/web"
	"github.com/dgrijalva/jwt-go"
)

var ErrInvalidToken = errors.New("Invalid token")

type Claims struct {
	UserID string   `json:"userid"`
	Roles  []string `json:"roles"`
	jwt.StandardClaims
}

func NewClaims(userid string, roles []string) Claims {
	c := Claims{
		UserID: userid,
		Roles:  roles,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}

	return c
}

func CreateToken(claims Claims) (string, error) {
	tkn := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)

	str, err := tkn.SignedString([]byte(os.Getenv("API_KEY")))
	if err != nil {
		return "", fmt.Errorf("generating token: %w", err)
	}

	return str, nil
}

func ParseWithClaims(tokenStr string) (Claims, error) {
	claims := Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("API_KEY")), nil
	})

	if err != nil {
		return Claims{}, web.NewRequestError(ErrInvalidToken, http.StatusBadRequest)
	}

	if !token.Valid {
		return Claims{}, web.NewRequestError(ErrInvalidToken, http.StatusBadRequest)
	}

	return claims, nil
}
