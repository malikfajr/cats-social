package config

import "github.com/golang-jwt/jwt/v5"

type CustomJWTClaim struct {
	ID    int              `json:"id"`
	Email string           `json:"email"`
	Name  string           `json:"name"`
	Exp   *jwt.NumericDate `json:"exp"`
	jwt.RegisteredClaims
}
