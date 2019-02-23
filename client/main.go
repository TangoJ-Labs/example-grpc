package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/seanbhart/example-grpc/protos"
	"google.golang.org/grpc"
)

const (
	address = "0.0.0.0:4000"
)

func main() {

	//=====================================
	//=============== SETUP ===============

	// Connect to the server
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("gRPC did not connect: %s\n", err)
	}
	defer conn.Close()

	//=============== SETUP ===============
	//=====================================

	// Call the rpc Iterate method and pass the Counter in a new Data object
	dataClient := pb.NewIterateCounterClient(conn)
	dataResponse, err := dataClient.Iterate(context.Background(), &pb.Data{Counter: 1})
	if err != nil {
		log.Fatalf("Error when calling Iterate: %s\n", err)
	}

	// Ouput the response from the server (the Counter value should have increased by 1)
	fmt.Printf("New counter value: %d\n", dataResponse.Counter)
}
