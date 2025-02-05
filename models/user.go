package models

import "gorm.io/gorm"

type User struct {
    gorm.Model
    Email          string `gorm:"uniqueIndex"`
    Password       string
    Role           string
    EmailVerified  bool
    VerificationToken string
    Name     string `json:"name"`
}