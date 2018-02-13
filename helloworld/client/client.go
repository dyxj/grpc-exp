package main

import (
	"context"
	"fmt"
	"log"
	"net"

	hw "github.com/dyxj/grpc-exp/helloworld/helloworld"
	"google.golang.org/grpc"
)

const (
	server     = "localhost"
	serverPort = "50051"
	clientName = "Mr. Client"
)

func main() {
	sAddr := net.JoinHostPort(server, serverPort)
	fmt.Println(sAddr)
	conn, err := grpc.Dial(sAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()
	// New client service
	c := hw.NewHelloWorldServiceClient(conn)
	r, err := c.SayHello(context.Background(), &hw.HelloRequest{Name: clientName})
	if err != nil {
		log.Printf("could not say hello: %v", err)
	}
	log.Printf("%s", r.Message)
	r, err = c.SayBye(context.Background(), &hw.HelloRequest{Name: clientName})
	if err != nil {
		log.Printf("could not say hello: %v", err)
	}
	log.Printf("%s", r.Message)
}
