package models

import "github.com/golang-jwt/jwt/v5"

type User struct {
	Username string
	Password string
	Role     string
}

type UserClaims struct {
	jwt.RegisteredClaims
	Username string
	Role     string
}
