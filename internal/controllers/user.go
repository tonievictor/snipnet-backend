package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"net/mail"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/siruspen/logrus"
	"golang.org/x/crypto/bcrypt"

	"snipnet/internal/utils"
	"snipnet/lib/services"
	"snipnet/lib/types"
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

func (u *UserController) UpdateUserOne(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(types.AuthSession).(types.Session)
	id := r.PathValue("id")

	var body types.UpdateOneData
	err := utils.ParseJson(r, &body)
	if err != nil {
		utils.WriteErr(w, http.StatusBadRequest, "No payload attached to req", err, u.log)
		return
	}

	if err = utils.Validate.Struct(body); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteErr(w, http.StatusBadRequest, "Missing parameters", error, u.log)
		return
	}

	usr, err := u.users.GetUser("id", id)
	if err != nil {
		utils.WriteErr(w, http.StatusNotFound, fmt.Sprintf("User with id %s not found", id), err, u.log)
		return
	}

	if usr.ID != session.UserID {
		utils.WriteErr(w, http.StatusUnauthorized, "You can't access this resource", errors.New("Unauthorized"), u.log)
		return
	}

	if body.Field != "email" && body.Field != "password" && body.Field != "username" {
		utils.WriteErr(w, http.StatusBadRequest, "You can't updated that parameter", errors.New("Invalid field Value"), u.log)
		return
	}

	if body.Field == "email" {
		_, err := mail.ParseAddress(body.Value)
		if err != nil {
			utils.WriteErr(w, http.StatusBadRequest, "Invalid email address", err, u.log)
			return
		}
	}

	if body.Field == "password" {
		hashedpwd, err := bcrypt.GenerateFromPassword([]byte(body.Value), bcrypt.DefaultCost)
		if err != nil {
			utils.WriteErr(w, http.StatusBadRequest, "An error occured while updating the resource", err, u.log)
			return
		}
		body.Value = string(hashedpwd)
	}

	user, err := u.users.UpdateUserSingle(id, body.Field, body.Value)
	if err != nil {
		utils.WriteErr(w, http.StatusBadRequest, "An error occured while updating the resource", err, u.log)
		return
	}

	utils.WriteRes(w, http.StatusOK, "User resource updated", user, u.log)
	return
}

func (u *UserController) UpdateUserMulti(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(types.AuthSession).(types.Session)
	id := r.PathValue("id")

	var body services.User
	err := utils.ParseJson(r, &body)
	if err != nil {
		utils.WriteErr(w, http.StatusBadRequest, "No payload attached to req", err, u.log)
		return
	}

	if err = utils.Validate.Struct(body); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteErr(w, http.StatusBadRequest, "Missing parameters", error, u.log)
		return
	}

	usr, err := u.users.GetUser("id", id)
	if err != nil {
		utils.WriteErr(w, http.StatusNotFound, fmt.Sprintf("User with id %s not found", id), err, u.log)
		return
	}

	if session.UserID != usr.ID {
		utils.WriteErr(w, http.StatusUnauthorized, "You can't access this resource", errors.New("Unauthorized"), u.log)
		return
	}

	user, err := u.users.UpdateUserMulti(&body)
	if err != nil {
		utils.WriteErr(w, http.StatusBadRequest, "An error occured while updating the resource", err, u.log)
		return
	}

	utils.WriteRes(w, http.StatusOK, "User resource updated", user, u.log)
	return
}

func (u *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(types.AuthSession).(types.Session)
	id := r.PathValue("id")

	usr, err := u.users.GetUser("id", id)
	if err != nil {
		utils.WriteErr(w, http.StatusNotFound, fmt.Sprintf("User with id %s not found", id), err, u.log)
		return
	}

	if session.UserID != usr.ID {
		utils.WriteErr(w, http.StatusUnauthorized, "You can't access this resource", errors.New("Unauthorized"), u.log)
		return
	}

	err = u.users.DeleteUser(id)
	if err != nil {
		utils.WriteErr(w, http.StatusBadRequest, "An error occured while updating the resource", err, u.log)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}
