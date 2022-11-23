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

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

var tools common.Tools

var failureCount = 0

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	log.Println("Processing authentication...")
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var tools common.Tools

	err := tools.ReadJSON(w, r, &requestPayload)
	if err != nil {
		_ = tools.ErrorJSON(w, err, http.StatusBadRequest)
		log.Panic(err)
		return
	}
	log.Printf("Receive request authenticate from user %s\n", requestPayload.Email)

	// check failure times
	if failureCount > 3 {
		mailPayload := MailPayload{
			From:    "admin@service.com",
			To:      requestPayload.Email,
			Subject: "Unusual log in detected",
			Message: "Your account has been log in unsuccessfully many times.",
		}

		app.sendMail(w, mailPayload)

		failureCount = 0
	}

	// validate the user against the database
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		_ = tools.ErrorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		log.Panic(err)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		failureCount++
		_ = tools.ErrorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		log.Panic(err)
		return
	}

	// log authentication
	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		_ = tools.ErrorJSON(w, err)
		log.Panic(err)
		return
	}

	payload := common.JSONResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	err = tools.WriteJSON(w, http.StatusAccepted, payload)
	if err != nil {
		log.Panic(err)
	}
}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.Marshal(entry)
	logServiceURL := fmt.Sprintf("http://logger-service:%s/log", common.LoggerPort)

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}
	return nil
}

func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	jsonData, _ := json.Marshal(msg)

	mailServiceURL := fmt.Sprintf("http://mailer-service:%s/send", common.MailerPort)

	// call the mail service
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		_ = tools.ErrorJSON(w, err)
		log.Panic(err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		_ = tools.ErrorJSON(w, err)
		log.Panic(err)
		return
	}
	defer response.Body.Close()

	// check back status code
	if response.StatusCode != http.StatusAccepted {
		_ = tools.ErrorJSON(w, errors.New("error calling mail service"))
		log.Panic(err)
		return
	}

	// send back json
	payload := common.JSONResponse{
		Error:   false,
		Message: "message sent to " + msg.To,
	}

	err = tools.WriteJSON(w, http.StatusAccepted, payload)
	if err != nil {
		log.Panic(err)
	}
}
