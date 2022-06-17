package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("secret_key") //TODO: add to config

type Token struct {
	Login string `json:"login,omitempty" db:"login"`
	jwt.StandardClaims
}

func Generate(data string) (string, time.Time, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	t := &Token{
		Login: data,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t)
	tokenString, err := token.SignedString(jwtKey)

	return tokenString, expirationTime, err
}

func Verify(tokenString string) (jwt.MapClaims, bool, bool) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)

	return claims, ok, token.Valid
}
