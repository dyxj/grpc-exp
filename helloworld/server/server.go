package main

import (
	"context"
	"log"
	"net"

	hw "github.com/dyxj/grpc-exp/helloworld/helloworld"
	"google.golang.org/grpc"
)

const (
	serverPort = ":50051"
)

// server is used to implement helloworld.HelloWorldServiceServer
type server struct{}

// SayHello: Returns hello message with input name
func (s *server) SayHello(ctx context.Context, in *hw.HelloRequest) (*hw.HelloResponse, error) {
	return &hw.HelloResponse{Message: "Hello World!!" + in.GetName()}, nil
}

// SayBye: Returns bye message with input name
func (s *server) SayBye(ctx context.Context, in *hw.HelloRequest) (*hw.HelloResponse, error) {
	return &hw.HelloResponse{Message: "Bye World!!" + in.GetName()}, nil
}

func main() {
	lis, err := net.Listen("tcp", serverPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	svr := grpc.NewServer()
	hw.RegisterHelloWorldServiceServer(svr, &server{})

	log.Println("starting Hello World rpc service on", serverPort)
	if err := svr.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
