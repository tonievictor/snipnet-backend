package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/siruspen/logrus"
	"golang.org/x/crypto/bcrypt"

	"snipnet/internal/utils"
	"snipnet/lib/services"
)

const SessionCookieName = "auth_token"

type SigninInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthController struct {
	users    services.UserStore
	sessions services.SessionStore
	log      *logrus.Logger
}

func NewAuthController(users services.UserStore, sessions services.SessionStore, log *logrus.Logger) *AuthController {
	return &AuthController{
		users:    users,
		sessions: sessions,
		log:      log,
	}
}

func (a *AuthController) Signup(w http.ResponseWriter, r *http.Request) {
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

	id := uuid.NewString()
	pwd, _ := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	hashedpwd := string(pwd)
	body.Password = hashedpwd
	body.ID = id

	user, err := a.users.CreateUser(&body)
	if err != nil {
		utils.WriteErr(w, http.StatusConflict, fmt.Sprintf("An error occured while creating user resource"), err, a.log)
		return
	}

	user.Password = ""
	utils.WriteRes(w, http.StatusCreated, "Account created successfully", user, a.log)
}

func (a *AuthController) Signin(w http.ResponseWriter, r *http.Request) {
	var body SigninInput

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
	auth := r.Header.Get("Authorization")
	if len(auth) != 0 && len(strings.Split(auth, " ")) != 1 {
		token := strings.Split(auth, " ")[1]
		if token != "" {
			sess, err := a.sessions.GetSession(token)
			if err == nil && time.Until(sess.ExpiryTime) > 0 {
				utils.WriteErr(w, http.StatusConflict, "Account is already logged in", errors.New("Account is logged in"), a.log)
				return
			}
		}
	}

	user, err := a.users.GetUser("username", body.Username)
	if err != nil {
		utils.WriteErr(w, http.StatusNotFound, "Invalid credentials", errors.New("Invalid credentials"), a.log)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		utils.WriteErr(w, http.StatusNotFound, "Invalid credentials", errors.New("Invalid credentials"), a.log)
		return
	}

	session_id := utils.GenerateSessionID()
	_, err = a.sessions.CreateSession(user.ID, session_id)
	if err != nil {
		utils.WriteErr(w, http.StatusInternalServerError, "An error occured while creating a new session", err, a.log)
		return
	}

	user.Password = ""
	user.AuthToken = session_id
	utils.WriteRes(w, http.StatusOK, "Account signed-in successfully", user, a.log)
}

// Sign out
func (a *AuthController) Signout(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if len(auth) == 0 || len(strings.Split(auth, " ")) != 2 {
		utils.WriteErr(w, http.StatusUnauthorized, "You are not logged in", errors.New("Session id not found"), a.log)
		return
	}

	token := strings.Split(auth, " ")[1]
	if token == "" {
		utils.WriteErr(w, http.StatusUnauthorized, "You are not logged in", errors.New("Session id not found"), a.log)
		return
	}

	_, err := a.sessions.GetSession(token)
	if err != nil {
		utils.WriteErr(w, http.StatusUnauthorized, "Invalid session token", err, a.log)
		return
	}

	err = a.sessions.DeleteSession(token)
	if err != nil {
		utils.WriteErr(w, http.StatusInternalServerError, "Unable to log you out, please try again", err, a.log)
		return
	}

	utils.WriteRes(w, http.StatusOK, "Account logged out successfully", "", a.log)
}
