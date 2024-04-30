package services

import (
	"context"
	"fmt"
	"time"
)

type UserStore interface {
	GetUser(field, value string) (*User, error)
	CheckUser(username, email string) (*User, error)
	CreateUser(user *User) (*User, error)
	DeleteUser(id string) error
	UpdateUserMulti(usr *User) (*User, error)
	UpdateUserSingle(id, field, value string) (*User, error)
}

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	AuthToken string    `json:"auth_token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) GetUser(field, value string) (*User, error) {
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := fmt.Sprintf("SELECT id, username, password, email, created_at, updated_at FROM users WHERE %s = $1;", field)

	row := db.QueryRowContext(ctx, query, value)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
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

	query := "SELECT id, username, password, email, created_at, updated_at FROM users WHERE username = $1 OR email = $2;"

	row := db.QueryRowContext(ctx, query, username, email)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) CreateUser(user *User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var saveduser User

	query := `
		INSERT INTO users (id, username, password, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, username, password, email, created_at, updated_at;
	`

	row := db.QueryRowContext(ctx, query, user.ID, user.Username, user.Password, user.Email, time.Now(), time.Now())
	err := row.Scan(
		&saveduser.ID,
		&saveduser.Username,
		&saveduser.Password,
		&saveduser.Email,
		&saveduser.CreatedAt,
		&saveduser.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &saveduser, nil
}

func (u *User) DeleteUser(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `DELETE FROM users WHERE id = $1;`

	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) GetUsers() (*[]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var users []*User

	query := "SELECT id, username, password, email, created_at, updated_at FROM users;"

	row, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		var user User
		err = row.Scan(
			&user.ID,
			&user.Username,
			&user.Password,
			&user.Email,
			user.CreatedAt,
			user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return &users, nil
}

func (u *User) UpdateUserSingle(id, field, value string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var user User

	query := fmt.Sprintf("UPDATE users SET %s = $1, updated_at = $2 WHERE id = $3 RETURNING id, username, email, created_at, updated_at", field)
	row := db.QueryRowContext(ctx, query, value, time.Now(), id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, err
}

func (u *User) UpdateUserMulti(usr *User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var user User

	query := "UPDATE users SET username = $1, email = $2,updated_at = $3 WHERE id = $4 RETURNING id, username, email, created_at, updated_at"
	row := db.QueryRowContext(ctx, query, usr.Username, usr.Email, time.Now(), usr.ID)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, err
}
