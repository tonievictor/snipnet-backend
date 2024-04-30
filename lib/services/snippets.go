package services

import (
	"context"
	"fmt"
	"time"

	"snipnet/lib/types"
)

type SnippetStore interface {
	GetSnippet(id string) (*types.SnippetWithUser, error)
	CreateSnippet(snippet *Snippet) (*Snippet, error)
	DeleteSnippet(id string) error
	UpdateSnippetMulti(snippet *Snippet) (*Snippet, error)
	UpdateSnippetSingle(id, field, value string) (*Snippet, error)
	GetSnippetsUser(user_id string) (*[]*types.SnippetWithUser, error)
	GetSnippets() (*[]*types.SnippetWithUser, error)
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

func (s *Snippet) GetSnippetsUser(user_id string) (*[]*types.SnippetWithUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var snippets []*types.SnippetWithUser

	query := `
		SELECT snippets.id, snippets.user_id, snippets.title, snippets.description, snippets.code, users.username, users.email, snippets.created_at, snippets.updated_at
		FROM snippets
		INNER JOIN users ON snippets.user_id = users.id
		WHERE snippets.user_id = $1;
	`
	row, err := db.QueryContext(ctx, query, user_id)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		var snippet types.SnippetWithUser

		err := row.Scan(
			&snippet.ID,
			&snippet.UserID,
			&snippet.Title,
			&snippet.Description,
			&snippet.Code,
			&snippet.Username,
			&snippet.Email,
			&snippet.CreatedAt,
			&snippet.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, &snippet)
	}

	return &snippets, nil
}

func (s *Snippet) GetSnippets() (*[]*types.SnippetWithUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var snippets []*types.SnippetWithUser

	query := `
		SELECT snippets.id, snippets.user_id, snippets.title, snippets.description, snippets.code, users.username, users.email, snippets.created_at, snippets.updated_at
		FROM snippets
		INNER JOIN users ON snippets.user_id = users.id
	`
	row, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		var snippet types.SnippetWithUser

		err := row.Scan(
			&snippet.ID,
			&snippet.UserID,
			&snippet.Title,
			&snippet.Description,
			&snippet.Code,
			&snippet.Username,
			&snippet.Email,
			&snippet.CreatedAt,
			&snippet.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, &snippet)
	}

	return &snippets, nil
}

func (s *Snippet) GetSnippet(id string) (*types.SnippetWithUser, error) {
	var snippet types.SnippetWithUser
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT snippets.id, snippets.user_id, snippets.title, snippets.description, snippets.code, users.username, users.email, snippets.created_at, snippets.updated_at
		FROM snippets
		INNER JOIN users ON snippets.user_id = users.id
		WHERE snippets.id = $1;
	`

	row := db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&snippet.ID,
		&snippet.UserID,
		&snippet.Title,
		&snippet.Description,
		&snippet.Code,
		&snippet.Username,
		&snippet.Email,
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
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, title, description, code, created_at, updated_at;
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

func (s *Snippet) DeleteSnippet(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := "DELETE from snippets WHERE id =$1;"
	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Snippet) UpdateSnippetSingle(id, field, value string) (*Snippet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var snip Snippet
	query := fmt.Sprintf(`
		UPDATE snippets
		SET %s = $1, updated_at = $2 WHERE id = $3
		RETURNING id, user_id, title, description, code, created_at, updated_at;
		`, field)

	row := db.QueryRowContext(ctx, query, value, time.Now(), id)
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

func (s *Snippet) UpdateSnippetMulti(snippet *Snippet) (*Snippet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var snip Snippet
	query := `
		UPDATE snippets
		SET title = $1, description = $2, code = $3, updated_at = $4
		WHERE id = $5 
		RETURNING id, user_id, title, description, code, created_at, updated_at;
	`
	row := db.QueryRowContext(ctx, query, snippet.Title, snippet.Description, snippet.Code, time.Now(), snippet.ID)
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
