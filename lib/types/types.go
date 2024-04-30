package types

import "time"

type Session struct {
	UserID     string
	SessionID  string
	CreatedAt  time.Time
	ExpiryTime time.Time
}

type SnippetWithUser struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Code        string    `json:"code" validate:"required"`
	Username    string    `json:"username" validate:"required"`
	Email       string    `json:"email" validate:"required,email"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

const AuthSession = "AuthSession"
