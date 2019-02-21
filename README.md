# gRPC Presentation - Simple w/ Authentication Example

## Create a self-signed cert & key for the server
- In the `server/auth` directory:
>`openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -keyout auth/key.pem -out auth/cert.pem`

## Compile the protocol buffer files into Go outputs
- In the `example` directory:
>`protoc -I protos/ -I ${GOPATH}/src --go_out=plugins=grpc:protos protos/example.proto`

<br>
<br>


## Run the Example
- In the `example` directory:
Server
>`go run server/main.go`

The server terminal should block and listen on the port

<br>

Client
>`go run client/main.go`

The client terminal should run the program (and complete), and output "New counter value: {input int + 1}"
The server terminal should output "Counter value: {input int}"