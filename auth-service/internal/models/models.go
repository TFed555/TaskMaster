package models

import (
	"time"
)

type Tokens struct {
	AccessToken  string  `db:"accesstoken"`
	RefreshToken string  `db:"refreshtoken"`
}

type User struct {
	ID        uint      `db:"id"`
	Login     string    `db:"login"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

type RefreshToken struct {
	ID		  uint 		`db:"id"`
	UserID    uint      `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
}
