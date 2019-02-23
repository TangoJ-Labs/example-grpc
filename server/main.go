package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/seanbhart/example-grpc/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	port = "localhost:4000"
)

type server struct{}

func (s *server) Iterate(ctx context.Context, data *pb.Data) (*pb.Data, error) {

	fmt.Printf("Counter value: %d\n", data.Counter)

	// Increase the Counter value by 1 and return
	data.Counter = data.Counter + 1
	return data, nil
}

// main starts a gRPC server and waits for connection
func main() {

	//=====================================
	//=============== SETUP ===============

	// Create a listener on TCP port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Error when creating Listener: %v\n", err)
	}

	// Create TransportCredentials
	creds, err := credentials.NewServerTLSFromFile("auth/cert.pem", "auth/key.pem")
	if err != nil {
		log.Fatalf("Error when creating TransportCredentials: %s\n", err)
	}

	// Create a gRPC server object
	grpcServer := grpc.NewServer(grpc.Creds(creds))

	// Create a server instance and register the server to receive messages
	dataServer := server{}
	pb.RegisterIterateCounterServer(grpcServer, &dataServer)

	// Start the gRPC server to listen on the port
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("Error when Serving: %s\n", err)
	}

	//=============== SETUP ===============
	//=====================================
}
