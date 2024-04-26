package services

import "time"

type Token struct {
	ID        int       `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}
