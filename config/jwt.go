package config

import "github.com/golang-jwt/jwt/v5"

type CustomJWTClaim struct {
	Email string           `json:"email"`
	Exp   *jwt.NumericDate `json:"exp"`
	jwt.RegisteredClaims
}
