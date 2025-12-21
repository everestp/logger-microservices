package main

import (
	"net/http"

	"github.com/everestp/log-service/data"
)


type JSONPayload struct{
	Name string `json:"name"`
	Data string `json:"data"`

}

func (app *Config) WriteLog(w http.ResponseWriter , r *http.Request){
    // read json in var
	var requestPayload JSONPayload
	app.readJSON(w, r , &requestPayload)
	// insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}
	err := app.Models.LogEntry.Insert(event)
	if err !=nil{
		app.errorJSON(w, err)
	}
	resp := jsonResponse{
		Error: false,
		Message: "logged",
	}
	app.writeJSON(w, http.StatusOK, resp)

}