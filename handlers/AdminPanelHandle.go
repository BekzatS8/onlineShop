package handlers

import (
    "encoding/json"
    "net/http"
    "onlinebooking/models"
    "gorm.io/gorm"
)

func AdminPanelHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "static/admin.htm")
    }
}

func ManageUsersHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var users []models.User
        if err := db.Find(&users).Error; err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(users)
    }
}

func ManageBookingsHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var bookings []models.Booking
        if err := db.Find(&bookings).Error; err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(bookings)
    }
}

func ManageCatalogHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var apartments []models.Apartment
        if err := db.Find(&apartments).Error; err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(apartments)
    }
}