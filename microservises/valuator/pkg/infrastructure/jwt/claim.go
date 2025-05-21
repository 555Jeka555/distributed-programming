package model

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	Login string `json:"login"`
	jwt.StandardClaims
}
