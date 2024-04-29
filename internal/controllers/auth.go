package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/siruspen/logrus"
	"golang.org/x/crypto/bcrypt"

	"snipnet/internal/utils"
	"snipnet/lib/services"
	"snipnet/lib/types"
)

const (
	SessionMaxAge = 259200
)

type SigninInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

var ctx = context.Background()

type AuthController struct {
	users services.UserStore
	log   *logrus.Logger
	cache *redis.Client
	ctx   context.Context
}

func NewAuthController(users services.UserStore, log *logrus.Logger, rds *redis.Client) *AuthController {
	return &AuthController{
		users: users,
		log:   log,
		cache: rds,
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
	return
}

func (a *AuthController) Signin(w http.ResponseWriter, r *http.Request) {
	var body SigninInput

	err := utils.ParseJson(r, &body)
	if err != nil {
		utils.WriteErr(w, http.StatusBadRequest, "No payload attached to req", err, a.log)
		return
	}

	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		token := strings.TrimPrefix(auth, "Bearer ")
		err := a.cache.Get(ctx, token).Err()
		if err != redis.Nil {
			_ = a.cache.Del(ctx, token).Err()
		}
	}

	if err = utils.Validate.Struct(body); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteErr(w, http.StatusBadRequest, "Missing parameters", error, a.log)
		return
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
	session, err := json.Marshal(types.Session{
		UserID:     user.ID,
		SessionID:  session_id,
		CreatedAt:  time.Now(),
		ExpiryTime: time.Now().Add(time.Hour * 24 * 3),
	})
	if err != nil {
		utils.WriteErr(w, http.StatusInternalServerError, "An error occured while creating a new session", err, a.log)
		return
	}
	err = a.cache.Set(ctx, session_id, session, time.Second*259200).Err()
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
	session := r.Context().Value(types.AuthSession).(types.Session)
	err := a.cache.Del(ctx, session.SessionID).Err()
	if err != nil {
		utils.WriteErr(w, http.StatusInternalServerError, "Unable to log you out, please try again", err, a.log)
		return
	}

	utils.WriteRes(w, http.StatusOK, "Account logged out successfully", "", a.log)
	return
}
