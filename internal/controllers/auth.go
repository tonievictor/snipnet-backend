package controllers

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/siruspen/logrus"
	"golang.org/x/crypto/bcrypt"

	"snipnet/internal/utils"
	"snipnet/lib/services"
)

type AuthController struct {
	users services.UserStore
	log   *logrus.Logger
}

func NewAuthController(users services.UserStore, log *logrus.Logger) *AuthController {
	return &AuthController{
		users: users,
		log:   log,
	}
}

func (a *AuthController) Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body services.User

	err := utils.ParseJson(r, &body)
	if err != nil {
		utils.WriteErr(w, http.StatusBadRequest, "No payload attached to req", err, a.log)
		return
	}

	if err = utils.Validate.Struct(body); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteErr(w, http.StatusBadRequest, "Missing parameters", error, a.log)
		return
	}

	u, err := a.users.CheckUser(body.Username, body.Email)
	if err == nil {
		var param string
		if u.Username == body.Username {
			param = body.Username
		} else {
			param = body.Email
		}
		utils.WriteErr(w, http.StatusConflict, fmt.Sprintf("%s already exists", param), fmt.Errorf("Account with %s already exists", param), a.log)
		return
	}

	id := uuid.New().String()
	pwd, _ := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	hashedpwd := string(pwd)
	body.Password = hashedpwd
	body.ID = id

	user, err := a.users.CreateUser(&body)
	if err != nil {
		utils.WriteErr(w, http.StatusConflict, fmt.Sprintf("An error occured while creating user resource"), err, a.log)
		return
	}

	utils.WriteRes(w, http.StatusCreated, "Account created successfully", user, a.log)
}
