package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"log/slog"

	"snipnet/lib/cache"
	"snipnet/lib/services"
)

func TestSignup(t *testing.T) {
	users := MockUser{}
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	rds := cache.Init()
	auth_ctrl := NewAuthController(&users, logger, rds)

	t.Run("should fail if an empty req.body is provided", func(t *testing.T) {
		payload := &MockUser{}

		body, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		auth_ctrl.Signup(rr, req)

		if rr.Result().StatusCode != http.StatusBadRequest {
			t.Fatal("Wrong status code")
		}
	})

	t.Run("should fail is some parameters are missing", func(t *testing.T) {
		payload := &services.User{
			Username: "testusrname",
			Email:    "test@gmail.com",
			// no password
		}
		body, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		auth_ctrl.Signup(rr, req)

		if rr.Result().StatusCode != http.StatusBadRequest {
			t.Fatal("Wrong status code")
		}
	})
}

// type UserStore interface {
// 	GetUsers() (*[]*User, error)
// 	GetUser(field, value string) (*User, error)
// 	CheckUser(username, email string) (*User, error)
// 	CreateUser(user *User) (*User, error)
// 	DeleteUser(id string) error
// 	UpdateUserMulti(usr *User) (*User, error)
// 	UpdateUserSingle(id, field, value string) (*User, error)

type MockUser struct{}

func (u *MockUser) GetUser(field, value string) (*services.User, error) {
	return nil, nil
}
 func (u *MockUser) GetUsers() (*[]*services.User, error)

func (u *MockUser) CheckUser(username, email string) (*services.User, error) {
	return nil, nil
}

func (u *MockUser) CreateUser(user *services.User) (*services.User, error) {
	return nil, nil
}

func (u *MockUser) UpdateUserMulti(user *services.User) (*services.User, error) {
	return nil, nil
}

func (u *MockUser) UpdateUserSingle(id, field, value string) (*services.User, error) {
	return nil, nil
}

func (u *MockUser) DeleteUser(id string) error {
	return nil
}
