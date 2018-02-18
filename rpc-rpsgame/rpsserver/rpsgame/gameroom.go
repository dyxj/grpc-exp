package rpsgame

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	winningScore = 3
)

type gameroommap struct {
	sync.RWMutex
	mapRooms map[string]gameroom
}

type gameroomslice struct {
	sync.RWMutex
	slRooms []*gameroom
}

type gameroom struct {
	RoomID    string
	Player1   *rpsSvcGameServer
	Player2   *rpsSvcGameServer
	IsFull    chan bool
	IsEnd     chan bool
	ScoreList []int
	// GameState
}

// JoinRoom : thread safe version for getRoom
func (grs *gameroomslice) JoinRoom(stream *rpsSvcGameServer) *gameroom {
	grs.Lock()
	defer grs.Unlock()
	return grs.joinRoom(stream)
}

// joinRoom : returns a gameroom from stack or creates one if there aren't any.
// When using JoinRoom, length of slice will always be only 0 or 1
// Plan to migrate to map later
func (grs *gameroomslice) joinRoom(stream *rpsSvcGameServer) *gameroom {
	var g *gameroom
	grsLen := len(grs.slRooms)
	if grsLen < 1 {
		// no rooms in slice
		// create and add to slice
		g = &gameroom{Player1: stream,
			IsFull:    make(chan bool),
			IsEnd:     make(chan bool),
			ScoreList: []int{0, 0},
		}
		grs.slRooms = append(grs.slRooms, g)
	} else {
		// get first
		g = grs.slRooms[0]
		// populate rooms
		if g.Player1 == nil {
			g.Player1 = stream
			return g
		} else if g.Player2 == nil {
			g.Player2 = stream
			// room is full remove from list
			close(g.IsFull)
			grs.slRooms = grs.slRooms[1:]
		}
	}
	return g
}

func (gr *gameroom) gameRoomMechanics() {
	rbegin := &Resp_Gstate{Gstate: Resp_BEGIN}
	resp := &Resp{Event: rbegin}
	gr.Player1.Send(resp)
	gr.Player2.Send(resp)
	// Channels to receive request
	// Get input from player 1 and 2
	for gr.ScoreList[0] < winningScore &&
		gr.ScoreList[1] < winningScore {
		// If context is canceled
		if gr.Player1.Context().Err() == context.Canceled ||
			gr.Player2.Context().Err() == context.Canceled {
			break
		}
		logrus.Infof("P1: %v\t\tP2: %v",
			gr.ScoreList[0], gr.ScoreList[1])
		// Send Get Input
		sendState(gr.Player1, Resp_ENTER_INPUT)
		sendState(gr.Player2, Resp_ENTER_INPUT)

		// Channel to receive signs
		p1SChan := make(chan Sign)
		p2SChan := make(chan Sign)
		go getPlayerSign(gr.Player1, p1SChan)
		go getPlayerSign(gr.Player2, p2SChan)
		p1Sign := <-p1SChan
		p2Sign := <-p2SChan

		// Process signs
		result := signLogic(p1Sign, p2Sign)
		if result == 1 {
			// player 1 win
			roundResults(gr.Player1, p2Sign, 0)
			roundResults(gr.Player2, p1Sign, 1)
			gr.ScoreList[0]++
		} else if result == 2 {
			// player 2 win
			roundResults(gr.Player1, p2Sign, 1)
			roundResults(gr.Player2, p1Sign, 0)
			gr.ScoreList[1]++
		} else {
			// Draw
			roundResults(gr.Player1, p2Sign, 2)
			roundResults(gr.Player2, p1Sign, 2)
		}
	}

	// Send Overall Win/Lose
	if gr.ScoreList[0] >= winningScore {
		matchResults(gr.Player1, true)
		matchResults(gr.Player2, false)
	} else {
		matchResults(gr.Player1, false)
		matchResults(gr.Player2, true)
	}
	close(gr.IsEnd)
}

func getPlayerSign(stream RpsSvc_GameServer, cSign chan Sign) {
	req, err := getRequest(stream)
	if err != nil {
		logrus.Errorf("Error at getPlayerSign(): %v, retry get player sign", err)
		state := &Resp_Gstate{Gstate: Resp_ERROR_REPEAT}
		resp := &Resp{Event: state}
		err = stream.Send(resp)
		if err != nil {
			logrus.Errorf("Fail to send error repeat at getPlayerSign(): %v, ending context", err)
			stream.Context().Done()
		}
		req, err = getRequest(stream)
		if err == nil {
			logrus.Errorf("Error at getPlayerSign() second trial: %v, defaulting to rock", err)
			cSign <- Sign_ROCK
			close(cSign)
			return
		}
	}

	if sign, ok := req.GetEvent().(*Req_Mysign); ok {
		cSign <- sign.Mysign
		close(cSign)
		return
	}
	logrus.Errorf("Not a valid sign, defaulting to rock")
	cSign <- Sign_ROCK
	close(cSign)
}

// 0 = draw
// 1 = p1
// 2 = p2
// if invalid sign
func signLogic(p1, p2 Sign) int {
	if p1 == Sign_PAPER {
		if p2 == Sign_ROCK {
			return 1
		} else if p2 == Sign_SCISSORS {
			return 2
		} else if p2 == Sign_PAPER {
			return 0
		}
	} else if p1 == Sign_ROCK {
		if p2 == Sign_ROCK {
			return 0
		} else if p2 == Sign_SCISSORS {
			return 1
		} else if p2 == Sign_PAPER {
			return 2
		}
	} else if p1 == Sign_SCISSORS {
		if p2 == Sign_ROCK {
			return 2
		} else if p2 == Sign_SCISSORS {
			return 0
		} else if p2 == Sign_PAPER {
			return 1
		}
	} else {
		// p1 invalid input
		// p2 valid input
		if p2 != Sign_SCISSORS && p2 != Sign_PAPER && p2 != Sign_ROCK {
			return 2
		}
		// both invalid input
		return 0
	}
	// p1 valid input
	// p2 invalid input
	return 1
}

// 0 win
// 1 lose
// 2 draw
func roundResults(stream RpsSvc_GameServer, oppSign Sign,
	wld int) {
	if wld == 0 {
		sendState(stream, Resp_WIN)
	} else if wld == 1 {
		sendState(stream, Resp_LOSE)
	} else if wld == 2 {
		sendState(stream, Resp_DRAW)
	}
	sendSign(stream, oppSign)
}

func matchResults(stream RpsSvc_GameServer, win bool) {
	if win {
		sendState(stream, Resp_OWIN)
	} else {
		sendState(stream, Resp_OLOSE)
	}
}
