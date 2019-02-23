package main

import (
	"context"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	pb "github.com/seanbhart/example-grpc/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	address = "localhost:4000"
)

func main() {

	//=====================================
	//=============== SETUP ===============

	// Read cert file
	pem, _ := ioutil.ReadFile("auth/cert.pem")

	// Create CertPool
	roots := x509.NewCertPool()
	roots.AppendCertsFromPEM(pem)

	// Create connection credentials
	creds := credentials.NewClientTLSFromCert(roots, "")

	// Connect to the server
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("ERROR: gRPC did not connect: %s\n", err)
	}
	defer conn.Close()

	//=============== SETUP ===============
	//=====================================

	//=====================================
	//== COMMAND LINE ARGUMENT HANDLING ===

	// Extract the counter starting argument (as a string)
	if len(os.Args) < 2 {
		log.Fatalln("Please include an integer argument as the counter starting point.")
	}
	arg := os.Args[1]

	// Covert the passed argument to int64
	intStart, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		log.Fatalf("That argument is not an integer.  Please try again.  Error: %s\n", err)
	}

	//== COMMAND LINE ARGUMENT HANDLING ===
	//=====================================

	// Call the rpc Iterate method and pass the Counter in a new Data object
	dataClient := pb.NewIterateCounterClient(conn)
	dataResponse, err := dataClient.Iterate(context.Background(), &pb.Data{Counter: intStart})
	if err != nil {
		log.Fatalf("Error when calling Iterate: %s\n", err)
	}

	// Ouput the response from the server (the Counter value should have increased by 1)
	fmt.Printf("New counter value: %d\n", dataResponse.Counter)
}
