package types

import "time"

type Session struct {
	UserID     string
	SessionID  string
	CreatedAt  time.Time
	ExpiryTime time.Time
}

const AuthSession = "AuthSession"
