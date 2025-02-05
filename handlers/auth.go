package handlers

import (
    "github.com/dgrijalva/jwt-go"
    "os"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
    Email string `json:"email"`
    Role  string `json:"role"`
    jwt.StandardClaims
}