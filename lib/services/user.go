package services

import (
	"context"
	"fmt"
	"time"

	"snipnet/lib/types"
)

type UserStore interface {
	GetUsers() (*[]*User, error)
	GetUser(field, value string) (*User, error)
	CheckUser(username, email string) (*User, error)
	CreateUser(id string, oauthUser *types.GHUser) (*User, error)
}

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Avatar    string    `json:"avatar"`
	AuthToken string    `json:"auth_token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) GetUser(field, value string) (*User, error) {
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := fmt.Sprintf("SELECT id, username, email, avatar, created_at, updated_at FROM users WHERE %s = $1;", field)

	row := db.QueryRowContext(ctx, query, value)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Avatar,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) CheckUser(username, email string) (*User, error) {
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := "SELECT id, username, email, avatar created_at, updated_at FROM users WHERE username = $1 OR email = $2;"

	row := db.QueryRowContext(ctx, query, username, email)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Avatar,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) CreateUser(id string, oauthUser *types.GHUser) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var saveduser User

	query := `
		INSERT INTO users (id, username, oauth_id, avatar, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, avatar, username, email, created_at, updated_at;
	`

	row := db.QueryRowContext(ctx, query, id, oauthUser.Name, string(oauthUser.ID), oauthUser.AvatarURL, oauthUser.Email, time.Now(), time.Now())
	err := row.Scan(
		&saveduser.ID,
		&saveduser.Avatar,
		&saveduser.Username,
		&saveduser.Email,
		&saveduser.CreatedAt,
		&saveduser.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &saveduser, nil
}

func (u *User) GetUsers() (*[]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var users []*User

	query := "SELECT id, username, email, avatar, created_at, updated_at FROM users;"

	row, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		var user User
		err = row.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Avatar,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return &users, nil
}
