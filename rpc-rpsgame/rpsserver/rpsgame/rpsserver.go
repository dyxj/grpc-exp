package rpsgame

import (
	"errors"

	"github.com/sirupsen/logrus"
)

// RpsSvc : implements RpsSvcServer interface for gRPC server
type RpsSvc struct {
	// map of games here
	rooms gameroomslice
}

// NewRpsSvc : creates and perform initialization
// of new Rps Service, returns it
func NewRpsSvc() (*RpsSvc, error) {
	rs := &RpsSvc{}
	return rs, nil
}

// Game : implementation of gRPC Game Service
func (r *RpsSvc) Game(stream RpsSvc_GameServer) error {
	logrus.Infof("%T", stream)
	// Get player from initial request
	req, err := getRequest(stream)
	if err != nil {
		return err
	}
	// Can get player id, will use it later
	player := req.GetJoin()
	// Ensure player is not nil
	if player == nil {
		logrus.Error("Player nil on initial join request")
		return errors.New("player was nil on initial join request")
	}
	logrus.Infof("Player attempting join %v", player)

	// Assign room
	streamcast := stream.(*rpsSvcGameServer)
	groom := r.rooms.JoinRoom(streamcast)

	// Wait for room to be full
	<-groom.IsFull
	if groom.Player2 == streamcast {
		logrus.Info("This is player two")
		groom.gameRoomMechanics()
		close(groom.IsEnd)
	}
	<-groom.IsEnd

	// Game Mechanics
	logrus.Info("Game End")
	return nil
}

func getRequest(stream RpsSvc_GameServer) (*Req, error) {
	req, err := stream.Recv()
	if err != nil {
		logrus.Errorf("getRequest(), Error receiving: %v", err)
		return nil, err
	}
	logrus.Infof("getRequest(), Received: %v", req)
	return req, nil
}

func sendResponse(stream RpsSvc_GameServer, r *Resp) error {
	logrus.Infof("sendResponse(), Sending response: %v", r)
	err := stream.Send(r)

	if err != nil {
		logrus.Errorf("sendResponse(), Error sending: %v", err)
	}

	return err
}

func sendSign(stream RpsSvc_GameServer, sign Sign) error {
	rsign := &Resp_Sign{Sign: sign}
	resp := &Resp{Event: rsign}
	return sendResponse(stream, resp)
}

func sendState(stream RpsSvc_GameServer, state Resp_State) error {
	rstate := &Resp_Gstate{Gstate: state}
	resp := &Resp{Event: rstate}
	return sendResponse(stream, resp)
}
