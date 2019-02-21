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
	// address = "docker.for.mac.host.internal:4000"
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
	// conn, err := grpc.Dial(address, grpc.WithInsecure())
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

	ctx := stream.Context()
	waitc := make(chan struct{})

	// Send random integer values and multiples to the server
	go func() {
		for i := 1; i <= 50; i++ {

			// Create a integer message object with random Value and Mutiple
			intMsg := pb.IntMsg{
				IntValue:    int64(rand.Intn(i)),
				IntMultiple: int64(rand.Intn(i)),
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
		err := stream.CloseSend()
		if err != nil {
			log.Println(err)
		}
	}()

	// Receive calculation responses
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

	// Wait for the context to be finished, then close the channel
	go func() {
		<-ctx.Done()
		err := ctx.Err()
		if err != nil {
			log.Println(err)
		}
		close(waitc)
	}()

	<-waitc
	log.Println("FINISHED")
}
