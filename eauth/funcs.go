package eauth

import (
	"crypto/rsa"
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

func ParserToken(tokenStr string, verifyKeys []*rsa.PublicKey) (claims AccountClaims, err error) {
	for _, verifyKey := range verifyKeys {
		token, err := jwt.ParseWithClaims(tokenStr, &AccountClaims{}, func(token *jwt.Token) (i interface{}, e error) {
			return verifyKey, nil
		})
		cla, ok := token.Claims.(*AccountClaims)
		if ok && token.Valid {
			return *cla, err
		}
	}
	return claims, errors.New("token invalid")
}
