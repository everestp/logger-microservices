package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/everestp/log-service/data"
	"github.com/everestp/log-service/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort ="8082"
	rpcPort ="5001"
	
)

var client *mongo.Client
var mongoUrl = env.GetString("MONGODB_URI", "")

type Config struct {
 Models data.Models
}

func main(){
	//connect to mongo
	mongoClient , err := connectToMongo()
	if err != nil{
		log.Panic(err)
	}
	client = mongoClient
	

	//create a context  in order to disconnect
	ctx ,cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()
	// close connection
	defer func()  {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
		
	}()
	
app := Config{
	Models: data.New(client),
}

// start web server
go app.server()



}

func (app *Config) server(){
	srv := &http.Server{
		Addr:fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
		
	}

	err :=srv.ListenAndServe()
	if err != nil{
		panic(err)
	}
}






func connectToMongo() (*mongo.Client , error){
	// create connection options
	clientOptions :=options.Client().ApplyURI(mongoUrl)

	//connect
	c ,err := mongo.Connect(context.TODO(),clientOptions)
	if err != nil {
		log.Println("Error Connecting to mongodb" ,err)
		return nil, err
	}
	return c,nil
}
