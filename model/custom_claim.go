package model

import "github.com/dgrijalva/jwt-go"

type MyCustomClaims struct {
	Account
	jwt.StandardClaims
}
