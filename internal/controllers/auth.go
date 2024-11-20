package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"

	"snipnet/internal/utils"
	"snipnet/lib/services"
	"snipnet/lib/types"
)

const (
	SessionMaxAge = 259200
)

var ctx = context.Background()

type AuthController struct {
	users       services.UserStore
	oauthConfig *oauth2.Config
	log         *slog.Logger
	cache       *redis.Client
}

func NewAuthController(
	users services.UserStore,
	oauthConfig *oauth2.Config,
	log *slog.Logger,
	rds *redis.Client,
) *AuthController {
	return &AuthController{
		users:       users,
		oauthConfig: oauthConfig,
		log:         log,
		cache:       rds,
	}
}

func (a *AuthController) GitHubOauth(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := a.oauthConfig.Exchange(r.Context(), code)
	if err != nil {
		utils.WriteErr(w, http.StatusBadRequest, "Invalid Code", err, a.log)
		return
	}

	ghUser, err := fetchGHUser(token.AccessToken)
	if err != nil {
		utils.WriteErr(w, http.StatusInternalServerError, "Invalid Code", err, a.log)
		return
	}

	user, err := a.users.GetUser("oauth_id", string(ghUser.ID))
	if err != nil {
		id := uuid.NewString()
		user, err = a.users.CreateUser(id, &ghUser)
		if err != nil {
			utils.WriteErr(
				w,
				http.StatusInternalServerError,
				"An error occured while creating the user resource",
				err, a.log
			)
			return
		}
	}

	session_id, err := a.createSession(user.ID)
	if err != nil {
		utils.WriteErr(
			w,
			http.StatusInternalServerError,
			"An error occured while creating a new session",
			err,
			a.log
		)
		return
	}

	user.AuthToken = session_id
	utils.WriteRes(w, http.StatusOK, "Account signed-in successfully", user, a.log)
	return
}

func (a *AuthController) createSession(userId string) (string, error) {
	session_id := utils.GenerateSessionID()
	session, err := json.Marshal(types.Session{
		UserID:     userId,
		SessionID:  session_id,
		CreatedAt:  time.Now(),
		ExpiryTime: time.Now().Add(time.Hour * 24 * 30),
	})
	if err != nil {
		return "", err
	}

	err = a.cache.Set(ctx, session_id, session, time.Second*2592000).Err()
	if err != nil {
		return "", err
	}

	return session_id, nil
}

func fetchGHUser(accessToken string) (types.GHUser, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return types.GHUser{}, nil
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return types.GHUser{}, nil
	}

	defer resp.Body.Close()
	var respBody types.GHUser
	if resp.Body == nil {
		return types.GHUser{}, err
	}

	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return types.GHUser{}, err
	}

	return respBody, nil
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
