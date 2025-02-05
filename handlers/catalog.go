package handlers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "onlinebooking/models"
    "strconv"

    "gorm.io/gorm"
)

func CatalogHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        page, _ := strconv.Atoi(r.URL.Query().Get("page"))
        if page < 1 {
            page = 1
        }
        pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
        if pageSize < 1 {
            pageSize = 10
        }
        sortBy := r.URL.Query().Get("sortBy")
        if sortBy == "" {
            sortBy = "created_at"
        }
        order := r.URL.Query().Get("order")
        if order == "" {
            order = "asc"
        }

        var apartments []models.Apartment
        offset := (page - 1) * pageSize
        if err := db.Order(fmt.Sprintf("%s %s", sortBy, order)).Offset(offset).Limit(pageSize).Find(&apartments).Error; err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(apartments); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    }
}

func CartHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Implement cart functionality here
    }
}
