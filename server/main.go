// See https://github.com/grpc/grpc-go for similar

package main

import (
	"io"
	"log"
	"net"
	"time"

	pb "github.com/seanbhart/example-grpc/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	port = "localhost:4000"
)

type server struct{}

func (s *server) Multiple(stream pb.BiDirectional_MultipleServer) error {

	// Keep the Mutiple server open to receive and send
	for {
		// Receive the incoming integer message object
		intMsg, err := stream.Recv()
		if err == io.EOF {
			// The server will receive an end of file error when the client is finished streaming data
			// Return nil to exit and close the server-side stream
			return nil
		}
		if err != nil {
			return err
		}

		// Multiply the Value integer by the Multiple integer and add the product to the object
		intCalc := intMsg.IntValue * intMsg.IntMultiple
		intMsg.IntCalc = intCalc

		// Delay to demonstrate the bidirectional streaming
		time.Sleep(time.Millisecond * 700)

		// Send the integer message object back to the client
		log.Printf("Sending response: %d * %d = %d", intMsg.IntValue, intMsg.IntMultiple, intMsg.IntCalc)
		err = stream.Send(intMsg)
		if err != nil {
			return err
		}
	}
}

// main starts a gRPC server and waits for connection
func main() {

	//=====================================
	//=============== SETUP ===============

	// Create a listener on TCP port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("ERROR when creating TransportCredentials: %v\n", err)
	}

	// Create TransportCredentials
	creds, err := credentials.NewServerTLSFromFile("auth/cert.pem", "auth/key.pem")
	if err != nil {
		log.Fatalf("ERROR: Failed to serve: %s\n", err)
	}

	// Create a gRPC server object
	grpcServer := grpc.NewServer(grpc.Creds(creds))

	// Create a server instance and register the server to receive messages
	biDirectionalServer := server{}
	pb.RegisterBiDirectionalServer(grpcServer, &biDirectionalServer)

	// Start the gRPC server to listen on the port
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("ERROR: Failed to serve: %s\n", err)
	}

	//=============== SETUP ===============
	//=====================================
}
