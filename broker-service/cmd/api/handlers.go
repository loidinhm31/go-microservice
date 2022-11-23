package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/loidinhm31/go-micro/common"
	"log"
	"net/http"
)

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

var tools common.Tools

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := common.JSONResponse{
		Error:   false,
		Message: "hit the broker",
	}

	var tools common.Tools

	err := tools.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		log.Println(err)
	}
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := tools.ReadJSON(w, r, &requestPayload)
	if err != nil {
		_ = tools.ErrorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.logItem(w, requestPayload.Log)
	default:
		_ = tools.ErrorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, _ := json.Marshal(entry)

	logServiceURL := fmt.Sprintf("http://logger-service:%s/log", common.LoggerPort)

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		_ = tools.ErrorJSON(w, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		_ = tools.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		_ = tools.ErrorJSON(w, err)
		return
	}

	payload := common.JSONResponse{
		Error:   false,
		Message: "logged",
	}
	err = tools.WriteJSON(w, http.StatusAccepted, payload)
	if err != nil {
		log.Println(err)
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// json to send to the auth microservice
	jsonData, _ := json.Marshal(a)

	authServiceURL := fmt.Sprintf("http://authentication-service:%s/authenticate", common.AuthPort)

	// call service
	request, err := http.NewRequest("POST", authServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		_ = tools.ErrorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		_ = tools.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// check back status code
	if response.StatusCode == http.StatusUnauthorized {
		_ = tools.ErrorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		_ = tools.ErrorJSON(w, errors.New("error calling auth service"))
		return
	}

	// read response body
	var jsonFromService common.JSONResponse

	// decode json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		_ = tools.ErrorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		_ = tools.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	payload := common.JSONResponse{
		Error:   false,
		Message: "authenticated",
		Data:    jsonFromService.Data,
	}

	err = tools.WriteJSON(w, http.StatusAccepted, payload)
	if err != nil {
		log.Println(err)
	}
}
