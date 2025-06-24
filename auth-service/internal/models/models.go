package models

import (
	"time"
)

type User struct {
	ID        uint      `db:"id"`
	Login     string    `db:"login"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

