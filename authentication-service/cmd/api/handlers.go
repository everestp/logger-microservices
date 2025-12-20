package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Autheticate(w  http.ResponseWriter , r *http.Request) {
	var requestPayload struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}
	err := app.ReadJSON(w, r , &requestPayload)
	if err != nil {
		app.errorJSON(w, err , http.StatusBadGateway)
		return
	}
  // validate the use against the database
  user ,err := app.Models.User.GetByEmail(requestPayload.Email)
  if err != nil{
	app.errorJSON(w, errors.New("invalid crenditals"), http.StatusBadRequest)
	return
  }
valid , err := user.PasswordMatches(requestPayload.Password)
if err != nil || !valid{
	app.errorJSON(w, errors.New("invalid crenditals"), http.StatusBadRequest)
	return
  }
payload :=jsonResponse{
	Error: false,
	Message: fmt.Sprintf("Logged in user %s",user.Email),
	Data: user,
}
app.WriteJSON(w, http.StatusAccepted, payload)
}