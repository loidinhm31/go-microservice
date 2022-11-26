package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/loidinhm31/go-microservice/common"
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

	// validate the user against the database
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		_ = tools.ErrorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		log.Println("Get username error:", err)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		failureCount++
		_ = tools.ErrorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		log.Println("Mismatch error:", err)

		// check failure times
		if failureCount > 3 {
			log.Printf("Failed %d time(s) from %s\n", failureCount, requestPayload.Email)

			mailPayload := MailPayload{
				From:    "admin@service.com",
				To:      requestPayload.Email,
				Subject: "Unusual log in detected",
				Message: "Your account has been log in unsuccessfully many times.",
			}

			err = app.sendMail(w, mailPayload)
			if err != nil {
				log.Println("Send email error:", err)
				return
			}
			failureCount = 0
		}
		return
	}

	// log authentication
	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		_ = tools.ErrorJSON(w, err)
		log.Println("Logger error:", err)
		return
	}

	payload := common.JSONResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	err = tools.WriteJSON(w, http.StatusAccepted, payload)
	if err != nil {
		log.Println("Write JSON error:", err)
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

func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) error {
	jsonData, _ := json.Marshal(msg)

	mailServiceURL := fmt.Sprintf("http://mailer-service:%s/send", common.MailerPort)

	// call the mail service
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// check back status code
	if response.StatusCode != http.StatusAccepted {
		return errors.New("error calling mail service")
	}

	// send back json
	payload := common.JSONResponse{
		Error:   false,
		Message: "message sent to " + msg.To,
	}

	err = tools.WriteJSON(w, http.StatusAccepted, payload)
	if err != nil {
		return err
	}
	return nil
}
