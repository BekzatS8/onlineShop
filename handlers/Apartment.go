package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"onlinebooking/models"
	"strconv"

	"gorm.io/gorm"
)

func AddApartmentHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            tmpl, err := template.ParseFiles("static/add_apartment.html")
            if err != nil {
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
            }
            tmpl.Execute(w, nil)
            return
        }

        if r.Method == http.MethodPost {
            r.ParseForm()
            name := r.FormValue("name")
            description := r.FormValue("description")
            price, err := strconv.ParseFloat(r.FormValue("price"), 64)
            if err != nil {
                http.Error(w, "Invalid price", http.StatusBadRequest)
                return
            }

            apartment := models.Apartment{
                Name:        name,
                Description: description,
                Price:       price,
            }

            if err := db.Create(&apartment).Error; err != nil {
                fmt.Println(err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
            }

            http.Redirect(w, r, "/index", http.StatusSeeOther)
        }
    }
}