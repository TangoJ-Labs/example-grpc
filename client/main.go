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

	// create a client connection
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("gRPC did not connect: %s\n", err)
	}
	defer conn.Close()

	// send the counter
	dataClient := pb.NewIterateCounterClient(conn)
	dataResponse, err := dataClient.Iterate(context.Background(), &pb.Data{Counter: 1})
	if err != nil {
		log.Fatalf("Error when calling IterateCounter: %s\n", err)
	}

	fmt.Printf("New counter value: %d\n", dataResponse.Counter)
}
