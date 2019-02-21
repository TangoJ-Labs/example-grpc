// Reference: https://github.com/grpc/grpc-go/blob/master/examples/route_guide/client/client.go

package main

import (
	"context"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	pb "github.com/seanbhart/example-grpc/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	address = "example-grpc-server:4000"
)

func main() {

	//=====================================
	//=============== SETUP ===============

	// Read cert file
	pem, _ := ioutil.ReadFile("auth/cert.pem")

	// Create CertPool
	roots := x509.NewCertPool()
	roots.AppendCertsFromPEM(pem)

	// create connection credentials
	creds := credentials.NewClientTLSFromCert(roots, "")

	// create a client connection
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("ERROR: gRPC did not connect: %s\n", err)
	}
	defer conn.Close()

	//=============== SETUP ===============
	//=====================================

	// create stream
	client := pb.NewBiDirectionalClient(conn)
	stream, err := client.Multiple(context.Background())
	if err != nil {
		log.Fatalf("ERROR opening sream: %v", err)
	}

	// Keep the client script running by creating a channel to block the script and wait until the stream is closed
	waitc := make(chan struct{})

	// Create a go routine to receive calculation responses
	go func() {
		for {
			intMsg, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("ERROR receiving: %v", err)
			}
			log.Printf("*** CALC: %d * %d = %d ***", intMsg.IntValue, intMsg.IntMultiple, intMsg.IntCalc)
		}
	}()

	// Send random integer values and multiples to the server
	for i := 1; i <= 30; i++ {

		// Create a integer message object with random Value and Mutiple
		intMsg := pb.IntMsg{
			IntValue:    int64(rand.Intn(20)),
			IntMultiple: int64(rand.Intn(20)),
		}

		// Send the object
		err := stream.Send(&intMsg)
		if err != nil {
			log.Fatalf("ERROR sending: %v", err)
		}
		log.Printf("Sent Value * Mutiple: %d * %d", intMsg.IntValue, intMsg.IntMultiple)

		// Delay to allow human viewing of processes
		time.Sleep(time.Millisecond * 200)
	}
	err = stream.CloseSend()
	if err != nil {
		log.Println(err)
	}

	// Client script channel to block until the grpc stream is closed
	<-waitc
	log.Println("FINISHED")
}
