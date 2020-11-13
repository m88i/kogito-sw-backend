package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type serverError struct {
	Message string `json:"message,omitempty"`
	Err     string `json:"error,omitempty"`
}

func (e serverError) Error() string {
	return e.Err
}

func newServerError(message string, err error) *serverError {
	return &serverError{Message: message, Err: err.Error()}
}

func handleServerError(serverError *serverError, writer http.ResponseWriter) {
	log.Error(serverError.Message+": ", serverError.Err)
	writer.WriteHeader(http.StatusInternalServerError)
	writer.Header().Add("Content-type", "application/json")
	if err := json.NewEncoder(writer).Encode(serverError); err != nil {
		log.Error("Failed to encode server error", err)
	}
}
