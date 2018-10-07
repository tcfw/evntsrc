package testinghelpers

import (
	jwt "github.com/dgrijalva/jwt-go"
)

//NulledJWT creates an empty JWT token for testing
func NulledJWT() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject: "000000000000000000000000",
	})

	ss, err := token.SignedString([]byte("TEST"))
	if err != nil {
		panic(err)
	}
	return ss
}
