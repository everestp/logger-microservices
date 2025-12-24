package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/everestp/log-service/data"
	"github.com/everestp/log-service/logs"
	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context , req *logs.LogRequest)(*logs.LogResponse , error){
input := req.GetLogEntry()


//write the log
logEntry :=data.LogEntry{
	Name: input.Name,
	Data: input.Data,
}
err :=l.Models.LogEntry.Insert(logEntry)
if err != nil {
	res := &logs.LogResponse{Result: "failed"}
	return res, err
}

//return response
res :=&logs.LogResponse{Result: "logged"}
return res, nil


}


func (app *Config) gRPCListen(){

	lis ,err := net.Listen("tcp", fmt.Sprintf(":%s", gRPCPort))
	if err !=nil{
		log.Fatalf("Failed to Listen for gRPC : %v",err)
	}
	s :=grpc.NewServer()
	logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models})
	log.Printf("gRPC Server started on port %s", gRPCPort)
	if err :=s.Serve(lis); err !=nil{
		log.Fatalf("Failed to Listen for gRPC : %v",err)
	}
	
}