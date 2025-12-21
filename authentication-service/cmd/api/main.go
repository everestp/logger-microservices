package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	// "os"
	
	"time"
_ "github.com/jackc/pgconn"
_ "github.com/jackc/pgx/v4"
_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/everestp/authentication/data"
)

const webPort ="8081"
var count int64
type Config struct {
 DB *sql.DB
 Models data.Models
}


func main() {
  log.Println("Starting authentication service")
  //TODO connect to database
 conn :=connectTODB()
 if conn == nil {
	log.Panic("Can't connect to Postgres!")
 }
  //set up config
  app :=Config{
	DB: conn,
	Models: data.New(conn),
  }

  srv := &http.Server{
	Addr: fmt.Sprintf(":%s",webPort),
	Handler: app.routes(),
	
  }
  err :=srv.ListenAndServe()
  if err !=nil{
	log.Panic(err)
  }
}

func openDB(dsn string) (*sql.DB , error){
	db ,err :=sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil

}
func connectTODB() *sql.DB {
	// dsn :=os.Getenv("DSN")
	for {
		connection ,err :=openDB("postgresql://neondb_owner:npg_cVwRfym2n8lb@ep-blue-base-adcmddes-pooler.c-2.us-east-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require")
		if err != nil{
  log.Println("Postgres not yet ready")
  count++
		} else{
			log.Print("Connect to Postgres!")
			return  connection
		}
		if count >10 {
			log.Println(err)
			return nil
		}
		log.Println("Backing off for two second")
		time.Sleep(2 *time.Second)
		 continue
	}
}