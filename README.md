# gRPC Presentation - Stream w/ Docker Example

## Create a self-signed cert & key for the server
In the repo root directory:
>`openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -keyout auth/key.pem -out auth/cert.pem`

## Compile the protocol buffer files into Go outputs
In the repo root directory:
>`protoc -I protos/ -I ${GOPATH}/src --go_out=plugins=grpc:protos protos/example.proto`

<br>
<br>


# SETUP - DOCKER NETWORK
Start the network separately to ensure the network name is consistent in all secondary docker-compose files.
<br>
<br>**0.1) Start the Docker Network**
>`docker network create example-grpc`

(you can check the running docker networks with `docker network list`)

## Run the Example
In the repo root directory:

Server
>`docker-compose -f server/docker-compose.yaml up -d`
>`docker exec -it example-grpc-server bash`
>`cd /go/src/github.com/seanbhart/example-grpc`
>`go run server/main.go`

The server terminal should block and listen on the port.

<br>

Client
>`docker-compose -f client/docker-compose.yaml up -d`
>`docker exec -it example-grpc-client bash`
>`cd /go/src/github.com/seanbhart/example-grpc`
>`go run client/main.go`

The client terminal will run the program and stream back responses from the server.
The server terminal will output some calculation logs.
The client terminal will display "FINISHED" when complete.