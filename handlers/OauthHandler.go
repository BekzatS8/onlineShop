package handlers

import (
    "net/http"
    "onlinebooking/models"
    "time"

    "github.com/dgrijalva/jwt-go"
    "github.com/markbates/goth/gothic"
    "gorm.io/gorm"
)

func OAuthCallbackHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        user, err := gothic.CompleteUserAuth(w, r)
        if err != nil {
            http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
            return
        }

        var dbUser models.User
        if err := db.Where("email = ?", user.Email).First(&dbUser).Error; err != nil {
            dbUser = models.User{
                Email: user.Email,
                Name:  user.Name,
            }
            db.Create(&dbUser)
        }

        expirationTime := time.Now().Add(24 * time.Hour)
        claims := &Claims{
            Email: dbUser.Email,
            StandardClaims: jwt.StandardClaims{
                ExpiresAt: expirationTime.Unix(),
            },
        }

        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
        tokenString, err := token.SignedString(jwtKey)
        if err != nil {
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