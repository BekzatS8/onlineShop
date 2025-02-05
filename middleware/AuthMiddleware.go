package middleware

import (
	"fmt"
	"net/http"
    "os"
	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
    Email string `json:"email"`
    jwt.StandardClaims
}

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("token")
        if err != nil {
            if err == http.ErrNoCookie {
                http.Redirect(w, r, "/login", http.StatusSeeOther)
                return
            }
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        tokenStr := cookie.Value
        claims := &Claims{}

        token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })
        if err != nil {
            fmt.Println(err)
            if err == jwt.ErrSignatureInvalid {
                http.Redirect(w, r, "/login", http.StatusSeeOther)
                return
            }
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
        if !token.Valid {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }

        next.ServeHTTP(w, r)
    })
}