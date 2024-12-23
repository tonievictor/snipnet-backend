package services

import (
	"context"
	"fmt"
	"time"

	"snipnet/types"
)

type SnippetStore interface {
	GetSnippet(id string) (*types.SnippetWithUser, error)
	CreateSnippet(snippet *Snippet) (*Snippet, error)
	DeleteSnippet(id string) error
	UpdateSnippetMulti(snippet *Snippet) (*Snippet, error)
	UpdateSnippetSingle(id, field, value string) (*Snippet, error)
	GetSnippetsUser(user_id string, offset, limit int, param, lang string) (*[]*types.SnippetWithUser, error)
	GetSnippets(offset, limit int, param, lang string) (*[]*types.SnippetWithUser, error)
}

type Snippet struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Language    string    `json:"language" validate:"required"`
	Code        string    `json:"code" validate:"required"`
	IsPublic    string    `json:"is_public" validate:"boolean"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (s *Snippet) GetSnippetsUser(user_id string, offset, limit int, param, lang string) (*[]*types.SnippetWithUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var snippets []*types.SnippetWithUser

	query := `
		SELECT snippets.id, snippets.user_id, snippets.title, snippets.description,
			snippets.language, snippets.code, snippets.is_public, users.username, users.email, users.avatar,
			snippets.created_at, snippets.updated_at
		FROM snippets
		INNER JOIN users ON snippets.user_id = users.id
		WHERE snippets.user_id = $1
			AND ($2 = '' OR document @@ to_tsquery($2))
			AND ($3 = '' OR snippets.language = $3)
			AND snippets.is_public
    ORDER BY snippets.updated_at DESC
		LIMIT $4
		OFFSET $5;
`
	row, err := db.QueryContext(ctx, query, user_id, param, lang, limit, offset)
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
			&snippet.Language,
			&snippet.Code,
			&snippet.IsPublic,
			&snippet.Username,
			&snippet.Email,
			&snippet.Avatar,
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

func (s *Snippet) GetSnippets(offset, limit int, param, lang string) (*[]*types.SnippetWithUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var snippets []*types.SnippetWithUser

	query := `
		SELECT snippets.id, snippets.user_id, snippets.title, snippets.description,
			snippets.language, snippets.code, snippets.is_public, users.username, users.email, users.avatar,
			snippets.created_at, snippets.updated_at
		FROM snippets
		INNER JOIN users ON snippets.user_id = users.id
		WHERE ($1 = '' OR document @@ to_tsquery($1))
			AND ($2 = '' OR snippets.language = $2)
			AND snippets.is_public
		ORDER BY snippets.updated_at DESC
		LIMIT $3
		OFFSET $4;
	`
	row, err := db.QueryContext(ctx, query, param, lang, limit, offset)
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
			&snippet.Language,
			&snippet.Code,
			&snippet.IsPublic,
			&snippet.Username,
			&snippet.Email,
			&snippet.Avatar,
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
		SELECT snippets.id, snippets.user_id, snippets.title, snippets.description,
			snippets.language, snippets.code, snippets.is_public, users.username, users.email, users.avatar,
			snippets.created_at, snippets.updated_at
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
		&snippet.Language,
		&snippet.Code,
		&snippet.IsPublic,
		&snippet.Username,
		&snippet.Email,
		&snippet.Avatar,
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
		INSERT INTO snippets (id, user_id, title, description, language ,code, is_public, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, user_id, title, description, language, code, is_public, created_at, updated_at;
	`

	row := db.QueryRowContext(ctx, query, snippet.ID, snippet.UserID,
		snippet.Title, snippet.Description, snippet.Language, snippet.Code, snippet.IsPublic,
		time.Now(), time.Now())
	err := row.Scan(
		&snip.ID,
		&snip.UserID,
		&snip.Title,
		&snip.Description,
		&snip.Language,
		&snip.Code,
		&snip.IsPublic,
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
		RETURNING id, user_id, title, description, language, code, is_public, created_at,
			updated_at;
		`, field)

	row := db.QueryRowContext(ctx, query, value, time.Now(), id)
	err := row.Scan(
		&snip.ID,
		&snip.UserID,
		&snip.Title,
		&snip.Description,
		&snip.Language,
		&snip.Code,
		&snip.IsPublic,
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
		SET title = $1, description = $2, language = $3, code = $4, updated_at = $5
		WHERE id = $6
		RETURNING id, user_id, title, description, code, is_public, created_at, updated_at;
	`
	row := db.QueryRowContext(ctx, query, snippet.Title, snippet.Description,
		snippet.Language, snippet.Code, time.Now(), snippet.ID)
	err := row.Scan(
		&snip.ID,
		&snip.UserID,
		&snip.Title,
		&snip.Description,
		&snip.Language,
		&snip.Code,
		&snip.IsPublic,
		&snip.CreatedAt,
		&snip.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &snip, nil
}
