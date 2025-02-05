package handlers

import (
    "net/http"
)

func LogoutHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        http.SetCookie(w, &http.Cookie{
            Name:   "token",
            Value:  "",
            MaxAge: -1,
        })

        http.Redirect(w, r, "/", http.StatusSeeOther)
    }
}