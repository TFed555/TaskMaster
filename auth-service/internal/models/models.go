package models

import (
	"time"
)

type User struct {
	ID uint `gorm:"primaryKey"`
	Login string
	Email string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null" json:"-"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

