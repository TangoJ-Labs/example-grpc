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

	// Create stream using the server's Multiple method
	client := pb.NewBiDirectionalClient(conn)
	stream, err := client.Multiple(context.Background())
	if err != nil {
		log.Fatalf("ERROR opening sream: %v", err)
	}

	// Keep the client script running by creating a channel to block the script and wait until the stream is closed.
	// We can use a simple channel of empty structures to use the channel for signalling rather than passing objects.
	// (Could use the WaitGroup (sync) package instead: https://golang.org/src/sync/example_test.go)
	waitc := make(chan struct{})

	// Create a go routine to receive calculation responses
	go func() {
		for {
			intMsg, err := stream.Recv()
			if err == io.EOF {
				// Look at the terminal output to see when the "EOF" error is sent. The server will send this error code
				// as an indicator that it has completed responding to the client and intends to shut down the stream.
				log.Println("EOF")

				// Now that the server has finished responding, the client does not need to keep blocking with the
				// signal channel, so close it down (drop down and finally pass line 113 below)
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("ERROR receiving: %v", err)
			}
			// Look for this output to indicate when the client has received back the server response for a request (calc in this case)
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

	// Look at the terminal output to see when the "CloseSend" command is sent from the client to the server.
	// The client will send this command to close the stream as soon as it is finished sending its data, but it will not
	// prematurely terminate the connection; the server will wait until it is finished responding to past requests to terminate the stream.
	log.Println("CloseSend")
	err = stream.CloseSend()
	if err != nil {
		log.Println(err)
	}

	// Client script channel to block until the grpc stream is closed
	<-waitc

	// The blocking "waitc" channel has been closed, output to show the script is ending
	log.Println("FINISHED")
}
