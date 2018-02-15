package main

import (
	"net"
	"os"

	"github.com/dyxj/grpc-exp/rpc-rpsgame/rpsserver/rpsgame"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	serverPort = "50051"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	lis, err := net.Listen("tcp", ":"+serverPort)
	if err != nil {
		logrus.Fatalf("Could not listen on port %v, %v", serverPort, err)
	}
	defer lis.Close()

	svr := grpc.NewServer()
	rpsSvc, err := rpsgame.NewRpsSvc()
	if err != nil {
		logrus.Fatalf("Could not create Rps service: %v", err)
	}

	rpsgame.RegisterRpsSvcServer(svr, rpsSvc)
	logrus.Infof("Starting server on port %v", serverPort)
	if err := svr.Serve(lis); err != nil {
		logrus.Fatalf("Server stopped: %v", err)
	}
}
