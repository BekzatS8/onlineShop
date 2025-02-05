package handlers

import (
	"fmt"
	"net/http"
	"onlinebooking/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Credentials struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func LoginHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var creds Credentials
        r.ParseForm()
        creds.Email = r.FormValue("email")
        creds.Password = r.FormValue("password")
        var user models.User
        if err := db.Where("email = ?", creds.Email).First(&user).Error; err != nil {
            fmt.Println(2)

            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
            fmt.Println(3)

            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        expirationTime := time.Now().Add(24 * time.Hour)
        claims := &Claims{
            Email: creds.Email,
            StandardClaims: jwt.StandardClaims{
                ExpiresAt: expirationTime.Unix(),
            },
        }

        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
        tokenString, err := token.SignedString(jwtKey)
        if err != nil {
            fmt.Println(4)

            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        http.SetCookie(w, &http.Cookie{
            Name:    "token",
            Value:   tokenString,
            Expires: expirationTime,
        })

        http.Redirect(w, r, "/index", http.StatusSeeOther)

    }
}