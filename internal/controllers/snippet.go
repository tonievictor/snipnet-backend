package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/siruspen/logrus"

	"snipnet/internal/utils"
	"snipnet/lib/services"
	"snipnet/lib/types"
)

type SnippetController struct {
	snippets services.SnippetStore
	log      *logrus.Logger
	cache    *redis.Client
}

func NewSnippetController(snippet services.SnippetStore, log *logrus.Logger, cache *redis.Client) *SnippetController {
	return &SnippetController{
		snippets: snippet,
		log:      log,
		cache:    cache,
	}
}

func (s *SnippetController) GetSnippetByID(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(types.AuthSession).(types.Session)
	id := r.PathValue("id")

	snippet, err := s.snippets.GetSnippet(id)
	if err != nil {
		utils.WriteErr(w, http.StatusNotFound, fmt.Sprintf("Snippet with %s not found", id), err, s.log)
		return
	}

	if session.UserID != snippet.UserID {
		utils.WriteErr(w, http.StatusUnauthorized, "You are not authorized to access this resource", errors.New("Not authorized"), s.log)
		return
	}

	utils.WriteRes(w, http.StatusOK, "Snippet found", snippet, s.log)
	return
}

func (s *SnippetController) CreateSnippet(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(types.AuthSession).(types.Session)
	var body services.Snippet

	err := utils.ParseJson(r, &body)
	if err != nil {
		utils.WriteErr(w, http.StatusBadRequest, "No payload attached to req", err, s.log)
		return
	}

	if err = utils.Validate.Struct(body); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteErr(w, http.StatusBadRequest, "Missing parameters", error, s.log)
		return
	}

	body.ID = uuid.NewString()
	body.UserID = session.UserID

	snippet, err := s.snippets.CreateSnippet(&body)
	if err != nil {
		utils.WriteErr(w, http.StatusInternalServerError, "An error occured while creating snippet", err, s.log)
		return
	}

	utils.WriteRes(w, http.StatusCreated, "Snippet created", snippet, s.log)
	return
}
