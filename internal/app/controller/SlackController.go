package controller

import (
	"fmt"
	"net/http"
)

type SlackController struct {
	Controller
}

func (controller SlackController) Entrypoint(response http.ResponseWriter, request *http.Request) {
	err := controller.HandleCORS(response, request)
	if err != nil {
		return
	}

	fmt.Printf("%#v\n", request)

	controller.Respond(response, "{}", http.StatusOK)
}
