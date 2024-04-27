package services

import (
	"context"
	"time"
)

type SessionStore interface {
	CreateSession(user_id, session_id string) (time.Duration, error)
	GetSession(session_id string) (*Session, error)
	DeleteSession(session_id string) error
}

type Session struct {
	ID         int       `json:"id"`
	UserID     string    `json:"user_id" validate:"required"`
	SessionID  string    `json:"session_id" validate:"required"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiryTime time.Time `json:"expiry_time" validate:"required"`
}

func (s *Session) CreateSession(user_id, session_id string) (time.Duration, error) {
	duration := time.Hour * 24 * 3 // 3 days
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		INSERT into sessions (session_id, user_id, created_at, expiry_time)
		VALUES ($1, $2, $3, $4);
	`

	_, err := db.ExecContext(ctx, query, session_id, user_id, time.Now(), time.Now().Add(duration))
	return duration, err
}

func (s *Session) GetSession(session_id string) (*Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var session Session

	query := `
		SELECT id, session_id, user_id, created_at, expiry_time FROM sessions
		WHERE session_id = $1;
	`

	row := db.QueryRowContext(ctx, query, session_id)
	err := row.Scan(
		&session.ID,
		&session.SessionID,
		&session.UserID,
		&session.CreatedAt,
		&session.ExpiryTime,
	)

	return &session, err
}

func (s *Session) DeleteSession(session_id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `DELETE FROM sessions WHERE session_id = $1;`

	_, err := db.ExecContext(ctx, query, session_id)
	return err
}
