package services

import (
	"context"
	"time"
)

type SnippetStore interface {
	GetSnippet(id string) (*Snippet, error)
	CreateSnippet(snippet *Snippet) (*Snippet, error)
}

type Snippet struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Code        string    `json:"code" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (s *Snippet) GetSnippet(id string) (*Snippet, error) {
	var snippet Snippet
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT id, user_id, title, description, code, created_at, updated_at FROM snippets WHERE id = $1;
	`

	row := db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&snippet.ID,
		&snippet.UserID,
		&snippet.Title,
		&snippet.Description,
		&snippet.Code,
		&snippet.CreatedAt,
		&snippet.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &snippet, nil
}

func (s *Snippet) CreateSnippet(snippet *Snippet) (*Snippet, error) {
	var snip Snippet

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		INSERT INTO snippets (id, user_id, title, description, code, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, user_id, title, description, code, created_at, updated_at;
	`

	row := db.QueryRowContext(ctx, query, snippet.ID, snippet.UserID, snippet.Title, snippet.Description, snippet.Code, time.Now(), time.Now())
	err := row.Scan(
		&snip.ID,
		&snip.UserID,
		&snip.Title,
		&snip.Description,
		&snip.Code,
		&snip.CreatedAt,
		&snip.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &snip, nil
}
