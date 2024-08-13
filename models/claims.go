package models

import "github.com/golang-jwt/jwt"

type Claims struct {
	ID             string             `json:"id"`
	Username       string             `json:"username"`
	Password       string             `json:"password"`
	Role           string             `json:"role"`
	StandardClaims jwt.StandardClaims `json:"standard_claims"`
}

func (c *Claims) Valid() error {
	return c.StandardClaims.Valid()
}
