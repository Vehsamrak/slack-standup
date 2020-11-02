package controller

import (
	"errors"
	"io"
	"net/http"
)

type Controller struct {
}

func (controller Controller) HandleCORS(response http.ResponseWriter, request *http.Request) error {
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Headers", "*")
	if request.Method == "OPTIONS" {
		return errors.New("HTTP OPTIONS requested")
	}

	return nil
}

// Respond sends HTTP response with headers
func (controller Controller) Respond(response http.ResponseWriter, body string, status int) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	io.WriteString(response, body)
}
