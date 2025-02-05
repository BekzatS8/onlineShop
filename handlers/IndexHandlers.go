package handlers

import (
    "html/template"
    "net/http"
    "onlinebooking/models"
    "gorm.io/gorm"
)

type IndexData struct {
    IsAuthenticated bool
    Apartments      []models.Apartment
}

func IndexHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        tmpl, err := template.ParseFiles("static/index.html")
        if err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        var apartments []models.Apartment
        if err := db.Find(&apartments).Error; err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        isAuthenticated := false
        if cookie, err := r.Cookie("token"); err == nil && cookie.Value != "" {
            isAuthenticated = true
        }

        data := IndexData{
            IsAuthenticated: isAuthenticated,
            Apartments:      apartments,
        }

        tmpl.Execute(w, data)
    }
}