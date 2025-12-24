package main

import (
	"context"
	"log"
	"time"

	"github.com/everestp/log-service/data"
)

type RPCServer struct {

}

type RPCPayload struct{
	Name string
	Data string
}

func (r *RPCServer)LogInfo(payload RPCPayload , resp *string) error {
	collection :=client.Database("logs").Collection("logs")
	_ , err :=collection.InsertOne(context.TODO(), data.LogEntry{
		Name: payload.Name,
		Data: payload.Data,
		CreatedAt: time.Now(),
	})

	if err != nil {
		log.Println("error writting to mongo",err)
		return err
	}
	*resp ="Processed payload via rpc" + payload.Name
	return nil
}