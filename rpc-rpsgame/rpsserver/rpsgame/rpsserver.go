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
	// Get player from initial request
	req, err := getRequest(stream)
	if err != nil {
		return err
	}
	player := req.GetJoin()
	// Ensure player is not nil
	if player == nil {
		logrus.Error("Player nil on initial join request")
		return errors.New("player was nil on initial join request")
	}
	logrus.Infof("Player attempting join %v", player)

	// Assign room
	r.rooms.GetRoom()

	// Check no. players in room
	// Wait for room to be full
	// Send Begin
	revent := &Resp_Gstate{Gstate: Resp_BEGIN}
	resp := &Resp{Event: revent}
	stream.Send(resp)

	reqNo := 1
	for {
		reqNo++
		logrus.Infof("Get Request %v", reqNo)
		req, err = getRequest(stream)
		if err != nil {
			return err
		}
		logrus.Infof("Request %v by %v", reqNo, player)
	}
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
