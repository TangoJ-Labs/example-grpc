package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/seanbhart/example-grpc/protos"
	"google.golang.org/grpc"
)

const (
	port = ":4000"
)

type server struct{}

func (s *server) Iterate(ctx context.Context, data *pb.Data) (*pb.Data, error) {

	fmt.Printf("Counter value: %d\n", data.Counter)

	// iterate the counter and return
	data.Counter = data.Counter + 1
	return data, nil
}

// main starts a gRPC server and waits for connection
func main() {
	// Create a listener on TCP port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	// Create a gRPC server object
	grpcServer := grpc.NewServer()

	// Create a server instance and register the server to receive messages
	dataServer := server{}
	pb.RegisterIterateCounterServer(grpcServer, &dataServer)

	// Start the gRPC listener
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %s\n", err)
	}
}
