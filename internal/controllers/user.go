package controllers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/redis/go-redis/v9"

	"snipnet/internal/utils"
	"snipnet/lib/services"
	"snipnet/lib/types"
)

type UserController struct {
	users services.UserStore
	log   *slog.Logger
	cache *redis.Client
}

func NewUserController(users services.UserStore, log *slog.Logger, rds *redis.Client) *UserController {
	return &UserController{
		users: users,
		log:   log,
		cache: rds,
	}
}

func (u *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := u.users.GetUsers()
	if err != nil {
		utils.WriteErr(w, http.StatusNotFound, "Unable to fetch users", err, u.log)
		return
	}

	utils.WriteRes(w, http.StatusOK, "Users found", users, u.log)
	return
}

func (u *UserController) GetUserByID(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(types.AuthSession).(types.Session)
	id := r.PathValue("id")

	user, err := u.users.GetUser("id", id)
	if err != nil {
		utils.WriteErr(w, http.StatusNotFound, fmt.Sprintf("User with id %s not found", id), err, u.log)
		return
	}

	if session.UserID != user.ID {
		utils.WriteErr(w, http.StatusUnauthorized, "You can't access this resource", errors.New("Unauthorized"), u.log)
		return
	}

	utils.WriteRes(w, http.StatusOK, "Users found", user, u.log)
	return
}
