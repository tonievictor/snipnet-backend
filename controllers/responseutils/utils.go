package responseutils

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func GenerateSessionID() string {
	return uuid.NewString()
}

func ParseJson(r *http.Request, payload interface{}) error {
	if r.Body == nil {
		return errors.New("Request payload missing")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}

var Validate = validator.New()

func WriteErr(w http.ResponseWriter, status_code int, message string, err error, log *slog.Logger) {
	log.Error("RESPONSE", slog.String("error", err.Error()))
	res := create_res(false, message, err.Error())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status_code)
	encode_err := json.NewEncoder(w).Encode(res)
	if encode_err != nil {
		log.Error("[RESPONSE]", slog.String("Error sending response to user:", encode_err.Error()))
	}
}

func WriteRes(w http.ResponseWriter, status_code int, message string, data interface{}, log *slog.Logger) {
	log.Info("RESPONSE", slog.String("message", message))
	res := create_res(true, message, data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status_code)
	encode_err := json.NewEncoder(w).Encode(res)
	if encode_err != nil {
		log.Error("[RESPONSE]", slog.String("Error sending response to user:", encode_err.Error()))
	}
}

func create_res(status bool, message string, data interface{}) Response {
	return Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
