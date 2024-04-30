package controllers

import (
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/siruspen/logrus"

	"snipnet/internal/utils"
	"snipnet/lib/services"
)

type UserController struct {
	users services.UserStore
	log   *logrus.Logger
	cache *redis.Client
}

func NewUserController(users services.UserStore, log *logrus.Logger, rds *redis.Client) *UserController {
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
