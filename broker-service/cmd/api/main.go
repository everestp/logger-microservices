package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)


const webPort ="8080"

type Config struct {
   RabbitMQ *amqp.Connection
}
func main(){
	//try to connect with rabbitmq
    rabbitConn , err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
   app := Config{
     RabbitMQ :rabbitConn,
   }

   
   //define http server
   srv :=&http.Server{
	   Addr: fmt.Sprintf(":%s",webPort),
	   Handler: app.routes() ,
	}
	
	// start the server
	err = srv.ListenAndServe()
	if err != nil{
		log.Panic(err)
	}
	log.Printf("Starting broker service on port %s\n",webPort)
}



//connect rabbitmq
func connect() (*amqp.Connection , error){
	var count int64
	var backoff = 1*time.Second
	var connection *amqp.Connection

	//do not continue until rabbitmq is  ready
	for {
		c , err :=amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready....")
			count++

		} else{
			connection =c
			 break
		}

		if count > 5{
			fmt.Println(err)
			return nil ,err
		}
		 backoff = time.Duration(math.Pow(float64(count),2))*time.Second
		 log.Println("backing off")
		 time.Sleep(backoff)
		 continue
	}
	return  connection , nil
}