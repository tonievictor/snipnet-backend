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

// @Summary      Get User
// @Description  Retrieve details of a specific user by their unique ID.
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "Unique ID of the user"
// @Success      200  {object}  services.User        "User details"
// @Failure      401  {object}  utils.Response       "Unauthorized access"
// @Failure      404  {object}  utils.Response       "User not found"
// @Router       /users/{id} [get]
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
