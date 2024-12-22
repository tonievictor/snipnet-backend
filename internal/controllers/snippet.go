package controllers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	validator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	redis "github.com/redis/go-redis/v9"

	"snipnet/internal/utils"
	"snipnet/lib/services"
	"snipnet/lib/types"
)

type SnippetController struct {
	snippets services.SnippetStore
	log      *slog.Logger
	cache    *redis.Client
}

func NewSnippetController(
	snippet services.SnippetStore,
	log *slog.Logger,
	cache *redis.Client,
) *SnippetController {
	return &SnippetController{
		snippets: snippet,
		log:      log,
		cache:    cache,
	}
}

/*
This function concatenates multiple req params together using ' & '
It trims white space around the string and gets rid of repeating spaces within the string
*/
func concatParam(s string) string {
	newStr := make([]byte, 0, len(s))
	sLen := len(s)
	for i := 0; i < sLen; i++ {
		if s[i] == ' ' {
			if i != (sLen-1) && s[i+1] != ' ' && len(newStr) > 0 {
				newStr = append(newStr, s[i])
				newStr = append(newStr, '&')
				newStr = append(newStr, s[i])
			}
			continue
		}
		newStr = append(newStr, s[i])

	}
	return string(newStr)
}

// @Summary      Delete Snippet
// @Description  Delete a snippet by its ID. Only the snippet owner can perform this action.
// @Tags         snippet
// @Security     ApiKeyAuth
// @Param        id    path     string  true  "Snippet ID to be deleted"
// @Success      204   "Snippet successfully deleted, no content returned"
// @Failure      401   {object} utils.Response  "Unauthorized access"
// @Failure      404   {object} utils.Response  "Snippet not found"
// @Failure      500   {object} utils.Response  "Internal server error during deletion"
// @Router       /snippets/{id} [delete]
func (s *SnippetController) DeleteSnippet(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(types.AuthSession).(types.Session)
	id := r.PathValue("id")

	snippet, err := s.snippets.GetSnippet(id)
	if err != nil {
		utils.WriteErr(w, http.StatusNotFound, fmt.Sprintf("Snippet with %s not found", id), err, s.log)
		return
	}

	if session.UserID != snippet.UserID {
		utils.WriteErr(w, http.StatusUnauthorized, `You are not authorized to access
			this resource`, errors.New("Not authorized"), s.log)
		return
	}

	err = s.snippets.DeleteSnippet(id)
	if err != nil {
		utils.WriteErr(
			w, http.StatusInternalServerError,
			"An error occured while deleting snippet",
			err,
			s.log,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}

// @Summary      Update Snippet Fields
// @Description  Update multiple fields of a snippet, such as title, description, and code.
// @Tags         snippet
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id    path     string              true  "Snippet ID to be updated"
// @Param        body  body     services.Snippet    true  "Updated snippet data"
// @Success      200	 {object} types.SnippetWithUser  "Updated Snippet details"
// @Failure      400   {object} utils.Response      "Invalid request or missing parameters"
// @Failure      401   {object} utils.Response      "Unauthorized access"
// @Failure      404   {object} utils.Response      "Snippet not found"
// @Failure      500   {object} utils.Response      "Internal server error during update"
// @Router       /snippets/{id} [patch]
func (s *SnippetController) UpdateSnippetMulti(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(types.AuthSession).(types.Session)
	id := r.PathValue("id")

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

	sp, err := s.snippets.GetSnippet(id)
	if err != nil {
		utils.WriteErr(w, http.StatusNotFound, fmt.Sprintf("Snippet with %s not found", id), err, s.log)
		return
	}

	if session.UserID != sp.UserID {
		utils.WriteErr(w, http.StatusUnauthorized, `You are not authorized to access
			this resource`, errors.New("Not authorized"), s.log)
		return
	}

	body.ID = sp.ID
	body.UserID = sp.UserID

	snippet, err := s.snippets.UpdateSnippetMulti(&body)
	if err != nil {
		utils.WriteErr(w, http.StatusInternalServerError, "Unable to update snippet", err, s.log)
		return
	}

	utils.WriteRes(w, http.StatusOK, "Updated snippet", snippet, s.log)
	return
}

// @Summary      Update Snippet
// @Description  Update a single field of a snippet, such as the title, description, or code.
// @Tags         snippet
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id    path     string              true  "Snippet ID to be updated"
// @Param        body  body     types.UpdateOneData true  "Field and value to update"
// @Success      200  {object} 	types.SnippetWithUser  "Updated Snippet details"
// @Failure      400   {object} utils.Response      "Invalid request or missing parameters"
// @Failure      401   {object} utils.Response      "Unauthorized access"
// @Failure      404   {object} utils.Response      "Snippet not found"
// @Router       /snippets/{id} [put]
func (s *SnippetController) UpdateSnippetOne(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(types.AuthSession).(types.Session)
	id := r.PathValue("id")

	var body types.UpdateOneData
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

	sp, err := s.snippets.GetSnippet(id)
	if err != nil {
		utils.WriteErr(
			w,
			http.StatusNotFound,
			fmt.Sprintf("Snippet with %s not found", id),
			err,
			s.log,
		)
		return
	}

	if session.UserID != sp.UserID {
		utils.WriteErr(
			w,
			http.StatusUnauthorized,
			"You are not authorized to access this resource",
			errors.New("Not authorized"),
			s.log,
		)
		return
	}

	if body.Field != "title" && body.Field != "description" && body.Field != "code" {
		utils.WriteErr(
			w,
			http.StatusBadRequest,
			"You can't update that parameter",
			errors.New("Invalid field Value"),
			s.log,
		)
		return
	}

	snippet, err := s.snippets.UpdateSnippetSingle(id, body.Field, body.Value)
	if err != nil {
		utils.WriteErr(
			w,
			http.StatusBadRequest,
			"An error occured while updating the resource",
			err,
			s.log,
		)
		return
	}

	utils.WriteRes(w, http.StatusOK, "Updated snippet", snippet, s.log)
	return
}

// @Summary      Get User's Snippets
// @Description  Retrieve all snippets created by a specific user, with optional filters.
// @Tags         snippet
// @Produce      json
// @Param        userid  path     string  true   "User ID whose snippets are being retrieved"
// @Param        page    query    string  false  "Page number for pagination (e.g., 1, 2, 3, ...)"
// @Param        param   query    string  false  "Search parameter to filter snippets"
// @Param        lang    query    string  false  "Programming language to filter snippets"
// @Success      200     {array} utils.Response  "List of snippets with user details"
// @Failure      404     {object} utils.Response  "Error fetching snippets"
// @Router       /users/{userid}/snippets [get]
func (s *SnippetController) GetAllUserSnippets(w http.ResponseWriter, r *http.Request) {
	user_id := r.PathValue("userid")
	query := r.URL.Query()
	page := query.Get("page")
	param := query.Get("param")
	param = concatParam(param)
	lang := query.Get("lang")
	var offset int
	limit := 20
	if p, err := strconv.Atoi(page); err != nil && p <= 0 {
		offset = 0
	} else {
		offset = (p - 1) * limit
	}
	snippets, err := s.snippets.GetSnippetsUser(user_id, offset, limit, param, lang)
	if err != nil {
		utils.WriteErr(w, http.StatusNotFound, "Error fetching snippets", err, s.log)
		return
	}
	utils.WriteRes(w, http.StatusOK, "User's snippets found", snippets, s.log)
	return
}

// @Summary      Get Snippets
// @Description  Retrieve all snippets, with optional filters.
// @Tags         snippet
// @Tags         snippet
// @Produce      json
// @Param        param  query     string  false  "Filter snippets by a specific string"
// @Param        page   query     string  false  "Page number, e.g., 0, 1, 2, ..."
// @Param        lang   query     string  false  "Programming language to filter snippets"
// @Success      200    {array}   types.SnippetWithUser  "List of snippets with user details"
// @Failure      500    {object}  utils.Response         "Internal server error"
// @Router       /snippets [get]
func (s *SnippetController) GetAllSnippets(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page := query.Get("page")
	param := query.Get("param")
	param = concatParam(param)
	lang := query.Get("lang")
	var offset int
	limit := 20
	if p, err := strconv.Atoi(page); err != nil && p <= 0 {
		offset = 0
	} else {
		offset = (p - 1) * limit
	}
	snippets, err := s.snippets.GetSnippets(offset, limit, param, lang)
	if err != nil {
		utils.WriteErr(w, http.StatusNotFound, "Error fetching snippets", err, s.log)
		return
	}

	if len(*snippets) == 0 {
		utils.WriteErr(w, http.StatusNotFound, "No Snippets found", errors.New(""), s.log)
		return
	}

	utils.WriteRes(w, http.StatusOK, "Snippets found", snippets, s.log)
	return
}

// @Summary      Get Snippet
// @Description  Retrieve a snippet by its unique ID.
// @Tags         snippet
// @Produce      json
// @Param        id   path     string  true  "Unique identifier for the snippet"
// @Success      200  {object} types.SnippetWithUser  "Snippet details along with user information"
// @Failure      500  {object} utils.Response         "Internal server error"
// @Router       /snippets/{id} [get]
func (s *SnippetController) GetSnippetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	snippet, err := s.snippets.GetSnippet(id)
	if err != nil {
		utils.WriteErr(w, http.StatusNotFound, fmt.Sprintf("Snippet with %s not found", id), err, s.log)
		return
	}

	utils.WriteRes(w, http.StatusOK, "Snippet found", snippet, s.log)
	return
}

// @Summary      Create Snippet
// @Description  Create a new snippet.
// @Tags         snippet
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        Authorization  header    string  true  "Bearer token for authentication"
// @Success      201  {object}  types.SnippetWithUser  "Created snippet details with user information"
// @Failure      500  {object}  utils.Response         "Internal server error"
// @Router       /snippets [post]
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
