# gRPC Presentation - Simple Example

## Compile the protocol buffer files into Go outputs
- From the "example" directory:
>`protoc -I protos/ -I ${GOPATH}/src --go_out=plugins=grpc:protos protos/example.proto`

<br>
<br>


## Run the Example
- From the "example" directory:
Server
>`go run server/main.go`

The server terminal should block and listen on the port

<br>

Client
>`go run client/main.go`

The client terminal should run the program (and complete), and output "New counter value: 2"
The server terminal should output "Counter value: 1"