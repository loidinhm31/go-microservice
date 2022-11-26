package main

import (
	"github.com/loidinhm31/go-microservice/common"
	"log"
	"net/http"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	log.Println("Processing mail before sending..")
	var tools common.Tools

	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage

	err := tools.ReadJSON(w, r, &requestPayload)
	if err != nil {
		_ = tools.ErrorJSON(w, err)
		log.Println("Read JSON error:", err)
		return
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	log.Printf("Execute send mail for email from %s\n", msg.From)
	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		_ = tools.ErrorJSON(w, err)
		log.Println("Send SMTP error:", err)
		return
	}

	payload := common.JSONResponse{
		Error:   false,
		Message: "sent to " + requestPayload.To,
	}

	err = tools.WriteJSON(w, http.StatusAccepted, payload)
	if err != nil {
		log.Println("Write JSON error:", err)
	}
	log.Printf("Mail has been sent to %s\n", msg.To)
}
