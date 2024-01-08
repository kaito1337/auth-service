package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Login        string `gorm:"uniqueIndex;not null" json:"login"`
	Email        string `gorm:"uniqueIndex;not null" json:"email"`
	Password     string `gorm:"not null" json:"-"`
	RefreshToken string `json:"refresh_token"`
}
