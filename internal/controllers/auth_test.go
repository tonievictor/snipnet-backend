package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/siruspen/logrus"

	"snipnet/lib/cache"
	"snipnet/lib/services"
)

func TestSignup(t *testing.T) {
	logger := logrus.Logger{
		Out: os.Stderr,
	}
	users := MockUser{}
	rds := cache.Init(os.Getenv("REDIS_CLIENT"))
	auth_ctrl := NewAuthController(&users, &logger, rds)

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

type MockUser struct{}

func (u *MockUser) GetUser(field, value string) (*services.User, error) {
	return nil, nil
}

func (u *MockUser) CheckUser(username, email string) (*services.User, error) {
	return nil, nil
}

func (u *MockUser) CreateUser(user *services.User) (*services.User, error) {
	return nil, nil
}

func (u *MockUser) UpdateUser(id, field, value string) (*services.User, error) {
	return nil, nil
}

func (u *MockUser) DeleteUser(id string) error {
	return nil
}
