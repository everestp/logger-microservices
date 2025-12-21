package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort ="8082"
	rpcPort ="5001"
	mongoUrl=""
)

var client *mongo.Client


type Config struct {

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
