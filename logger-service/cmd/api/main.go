
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
    webPort = "8082"
    rpcPort = "5001"
)

var client *mongo.Client
var mongoUrl = env.GetString("MONGODB_URI", "mongodb+srv://everest:kaiseho12@cluster0.hodtqt5.mongodb.net/?appName=Cluster0")

type Config struct {
    Models data.Models
}

func main() {
    fmt.Println("MONGODB_URI", mongoUrl)
    
    // connect to mongo
    mongoClient, err := connectToMongo()
    if err != nil {
        log.Panic(err)
    }
    client = mongoClient

    // create a context in order to disconnect
    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()
    
    // close connection
    defer func() {
        if err = client.Disconnect(ctx); err != nil {
            panic(err)
        }
    }()

    app := Config{
        Models: data.New(client),
    }

    // start web server
    // Removed "go" keyword so the main process waits for the server to finish
    app.server() 
}

func (app *Config) server() {
    srv := &http.Server{
        Addr:    fmt.Sprintf(":%s", webPort),
        Handler: app.routes(),
    }

	log.Printf("Starting logger service on port %s\n",webPort)
    err := srv.ListenAndServe()
    if err != nil {
        panic(err)
    }
}

func connectToMongo() (*mongo.Client, error) {
    // create connection options
    clientOptions := options.Client().ApplyURI(mongoUrl)

    // connect
    c, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Println("Error Connecting to mongodb:", err)
        return nil, err
    }

    // Verify the connection is actually successful
    err = c.Ping(context.TODO(), nil)
    if err != nil {
        log.Println("Database reachable check failed:", err)
        return nil, err
    }

    log.Println("Connected to MongoDB successfully!")
    return c, nil
}

