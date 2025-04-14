package models

import "time"

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // не отдаем пароль в API
	CreatedAt time.Time `json:"created_at"`
}
