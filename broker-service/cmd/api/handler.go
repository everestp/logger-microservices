package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"time"

	"net/http"
	"net/rpc"

	"github.com/everestp/broker-service/event"
	"github.com/everestp/broker-service/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// RequestPayload defines the incoming JSON structure from the frontend
type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}
	_ = app.writeJSON(w, http.StatusOK, payload)
}



func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.logITemViaRPC(w, requestPayload.Log)
	case "mail":
		app.SendMail(w, requestPayload.Mail)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) SendMail(w http.ResponseWriter, msg MailPayload) {
	jsonData, _ := json.Marshal(msg)

	// URL matches the service name and internal port in docker-compose
	mailServiceURL := "http://mailer-service:8083/send"

	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling mailer service"))
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Message sent to " + msg.To
	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, _ := json.Marshal(a)
	url := "http://authentication-service:8081/authenticate"

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, _ := json.Marshal(entry)
	url := "http://logger-service:8082/log"

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("log service failed"))
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"
	app.writeJSON(w, http.StatusAccepted, payload)
}


func (app *Config) logEventViaRabbit(w http.ResponseWriter , l LogPayload){
  err := app.pushTOQueue(l.Name,l.Data)
if err !=nil{
	app.errorJSON(w, err)
	return
}
var payload jsonResponse
payload.Error= false
payload.Message ="logged via rabbitmq"
app.writeJSON(w, http.StatusAccepted, payload)
}


//pushToQueue pushes a message to RabbitMQ
func (app *Config) pushTOQueue(name , msg string) error {

	emitter , err :=event.NewEventEmmiter(app.RabbitMQ)
	if err != nil {
		return err
	}
	payload :=LogPayload{
	Name: name,
	Data: msg,
	}

	j ,_ :=json.MarshalIndent(&payload, "", "\t")
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil

}
type RPCPayload struct {
	Name string
	Data string
}

func(app *Config) logITemViaRPC(w http.ResponseWriter , l LogPayload) {
	client , err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil{
		app.errorJSON(w, err )
		return
	}
	rpcPayload :=RPCPayload(l)
	var result string
	err = client.Call("RPCServer.LogInfo",rpcPayload, &result)
	if err != nil{
		app.errorJSON(w, err)
     return
	}
	paylod :=jsonResponse{
		Error: false,
		Message: result,
	}
	app.writeJSON(w, http.StatusAccepted, paylod)
}


func (app *Config) LogViaGRPC(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	conn, err := grpc.Dial("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer conn.Close()

	c := logs.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log {
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"

	app.writeJSON(w, http.StatusAccepted, payload)
}