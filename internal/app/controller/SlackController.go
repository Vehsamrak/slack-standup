package controller

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vehsamrak/slack-standup/internal/app/controller/requests"
	"io/ioutil"
	"net/http"
)

type SlackController struct {
	Controller
}

func (controller SlackController) Ping(response http.ResponseWriter, request *http.Request) {
	err := controller.HandleCORS(response, request)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()

	ping := &requests.Ping{}
	err = json.Unmarshal(body, &ping)
	if err != nil {
		panic(err)
	}

	log.Infof("Ping requested: %#v", string(body))

	controller.Respond(response, fmt.Sprintf("{\"challenge\":\"%s\"}", ping.Challenge), http.StatusOK)
}

func (controller SlackController) Entrypoint(response http.ResponseWriter, request *http.Request) {
	err := controller.HandleCORS(response, request)
	if err != nil {
		return
	}

	fmt.Printf("%#v\n", request)

	controller.Respond(response, "{}", http.StatusOK)
}
