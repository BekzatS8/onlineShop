package models

import "gorm.io/gorm"

type Apartment struct {
    gorm.Model
    Name        string
    Description string
    Price       float64
}