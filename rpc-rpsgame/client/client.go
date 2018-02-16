package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/dyxj/grpc-exp/rpc-rpsgame/rpsserver/rpsgame"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	server     = "localhost"
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
	svrAddr := net.JoinHostPort(server, serverPort)
	logrus.Info("Connecting to" + svrAddr)
	conn, err := grpc.Dial(svrAddr, grpc.WithInsecure())
	if err != nil {
		logrus.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()
	// New client service
	c := rpsgame.NewRpsSvcClient(conn)

	// Calls game service
	gclient, err := c.Game(context.Background())
	if err != nil {
		logrus.Fatalf("error accessing game: %v", err)
	}

	// Sends join request
	player := &rpsgame.Req_Player{Id: "1234"}
	rj := &rpsgame.Req_Join{Join: player}
	rq := &rpsgame.Req{Event: rj}
	err = gclient.Send(rq)
	if err != nil {
		logrus.Errorf("error from game server: %v", err)
	}

	// Gets game state, begin
	resp, err := gclient.Recv()
	if err != nil {
		logrus.Errorf("error receiving: %v", err)
	}
	evt := resp.GetEvent()
	logrus.Info(evt)

	// ctx := gclient.Context()
	// fmt.Println(ctx)
	// Add on this part later
	cont := true
	if cont {
		// Handle user input
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("Enter input: ")
			text, err := reader.ReadString('\n')
			if err != nil {
				logrus.Error("Error getting input")
				continue
			}
			switch strings.Trim(text, "\n") {
			case "q", "quit", "exit":
				gclient.CloseSend()
				return
			default:
				sendSign(gclient, rpsgame.Sign_PAPER)
				fmt.Println(text)
			}

			resp, err := gclient.Recv()
			if err != nil {
				logrus.Errorf("error receiving: %v", err)
			}
			evt := resp.GetEvent()
			logrus.Info(evt)
		}
	}

}

func sendSign(stream rpsgame.RpsSvc_GameClient, sign rpsgame.Sign) error {
	rsign := &rpsgame.Req_Mysign{Mysign: sign}
	rq := &rpsgame.Req{Event: rsign}
	return sendRequest(stream, rq)
}

func sendRequest(stream rpsgame.RpsSvc_GameClient, req *rpsgame.Req) error {
	logrus.Infof("sendRequest(), Sending response: %v", req)
	err := stream.Send(req)
	if err != nil {
		logrus.Errorf("sendResponse(), Error sending: %v", err)
	}

	return err
}
