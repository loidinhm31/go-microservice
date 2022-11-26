package main

import (
	"github.com/loidinhm31/go-microservice/common"
	"log"
	"log-service/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

var tools common.Tools

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload JSONPayload
	_ = tools.ReadJSON(w, r, &requestPayload)

	// insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		_ = tools.ErrorJSON(w, err)
		log.Println("Write event error:", err)
		return
	}

	response := common.JSONResponse{
		Error:   false,
		Message: "logged",
	}

	err = tools.WriteJSON(w, http.StatusAccepted, response)
	if err != nil {
		log.Println("Write JSON error:", err)
	}
}
