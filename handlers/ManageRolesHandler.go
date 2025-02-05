package handlers

import (
    "encoding/json"
    "net/http"
    "onlinebooking/models"
    "gorm.io/gorm"
    "log"
)

func ManageRolesHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            var users []models.User
            if err := db.Find(&users).Error; err != nil {
                log.Printf("Error retrieving users: %v", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
            }

            w.Header().Set("Content-Type", "application/json")
            if err := json.NewEncoder(w).Encode(users); err != nil {
                log.Printf("Error encoding users to JSON: %v", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
            }
        }

        if r.Method == http.MethodPost {
            email := r.FormValue("email")
            role := r.FormValue("role")

            var user models.User
            if err := db.Where("email = ?", email).First(&user).Error; err != nil {
                log.Printf("User not found: %v", err)
                http.Error(w, "User not found", http.StatusNotFound)
                return
            }

            user.Role = role
            if err := db.Save(&user).Error; err != nil {
                log.Printf("Error updating user role: %v", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
            }

            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(map[string]string{"message": "Role updated successfully"})
        }
    }
}