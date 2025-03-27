package models

import "github.com/golang-jwt/jwt/v4"

type TokenClaims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

type Oauth struct {
	Code string `json:"code"`
}
