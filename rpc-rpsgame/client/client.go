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
	logrus.SetLevel(logrus.ErrorLevel)
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

	// Enter player id
	fmt.Print("Enter player id: ")
	reader := bufio.NewReader(os.Stdin)
	playerid, err := reader.ReadString('\n')
	if err != nil {
		logrus.Fatal("Error getting player id")
	}

	// Sends join request
	player := &rpsgame.Req_Player{Id: playerid}
	rj := &rpsgame.Req_Join{Join: player}
	rq := &rpsgame.Req{Event: rj}
	err = gclient.Send(rq)
	if err != nil {
		logrus.Errorf("error from game server: %v", err)
	}

	// Gets game state, begin
	fmt.Println("Waiting for opponent...")
	resp, err := gclient.Recv()
	if err != nil {
		logrus.Errorf("error receiving: %v", err)
	}
	evt := resp.GetEvent()
	resp.GetGstate()

	state, ok := resp.GetEvent().(*rpsgame.Resp_Gstate)
	if !ok {
		logrus.Fatal("Expecting Game State...")
	}
	if state.Gstate != rpsgame.Resp_BEGIN {
		logrus.Info(evt)
		logrus.Fatal("Expecting State Begin...")
	}
	fmt.Println("Game start")
	fmt.Println("You can quit any time by pressing ctrl+c")
	gameState := state.Gstate
	lastMySign := rpsgame.Sign_ROCK
	ScoreList := []int{0, 0}

	// Game loop
	for gameState != rpsgame.Resp_OWIN && gameState != rpsgame.Resp_OLOSE {
		// Get event
		evt := getEvent(gclient)
		// Update Game State or Print sign and continue
		if evtSign, ok := evt.(*rpsgame.Resp_Sign); ok {
			fmt.Printf("My Sign\t\tDude's Sign\n%v\t\t%v\n",
				lastMySign, evtSign.Sign)
			fmt.Println("-----------------------------")
			continue
		} else if evtState, ok := evt.(*rpsgame.Resp_Gstate); ok {
			gameState = evtState.Gstate
		}

		// Handle game states
		if gameState == rpsgame.Resp_ENTER_INPUT {
			// Input
			fmt.Print("Rock(1), Paper(2), Scissors(3), Quit(q)\nEnter input: ")
			text, err := reader.ReadString('\n')
			if err != nil {
				logrus.Error("Error getting input")
				continue
			}
			switch strings.Trim(text, "\n") {
			case "q", "quit", "exit":
				gclient.CloseSend()
				// return
				break
			case "1":
				sendSign(gclient, rpsgame.Sign_ROCK)
				lastMySign = rpsgame.Sign_ROCK
			case "2":
				sendSign(gclient, rpsgame.Sign_PAPER)
				lastMySign = rpsgame.Sign_PAPER
			case "3":
				sendSign(gclient, rpsgame.Sign_SCISSORS)
				lastMySign = rpsgame.Sign_SCISSORS
			default:
				fmt.Println("Invalid input!")
				continue
			}
			fmt.Println("Waiting for the dude")
		} else if gameState == rpsgame.Resp_WIN {
			ScoreList[0]++
			fmt.Println("You won this round")
			fmt.Printf("Score\nYou\t\tDude's Score\n%v\t\t%v\n", ScoreList[0], ScoreList[1])
			continue
		} else if gameState == rpsgame.Resp_LOSE {
			ScoreList[1]++
			fmt.Println("You lost this round")
			fmt.Printf("Score\nYou\t\tDude's Score\n%v\t\t%v\n", ScoreList[0], ScoreList[1])
			continue
		} else if gameState == rpsgame.Resp_DRAW {
			fmt.Println("It's a draw this round")
			continue
		} else if gameState == rpsgame.Resp_OWIN {
			fmt.Println("You've won the match!!!")
			gclient.CloseSend()
			break
		} else if gameState == rpsgame.Resp_OLOSE {
			fmt.Println("You've lost the match...")
			gclient.CloseSend()
			break
		} else if gameState == rpsgame.Resp_ERROR_REPEAT {
			fmt.Println("An error occured, repeat game state")
			continue
		}
		// Handle game states end
		fmt.Println("-----------------------------")
	}
	// Game Loop End
	fmt.Println("Game closing")
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

func getEvent(gclient rpsgame.RpsSvc_GameClient) interface{} {
	resp, err := gclient.Recv()
	if err != nil {
		logrus.Errorf("error receiving: %v", err)
	}
	return resp.GetEvent()
}
