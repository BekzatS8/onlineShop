package models

import "gorm.io/gorm"

type Booking struct {
    gorm.Model
    UserID uint
    Details string
}