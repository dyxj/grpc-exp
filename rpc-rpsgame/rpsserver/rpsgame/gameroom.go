package rpsgame

import (
	"sync"
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
	RoomID  string
	Player1 *rpsSvcGameServer
	Player2 *rpsSvcGameServer
	IsFull  chan bool
}

// JoinRoom : thread safe version for getRoom
func (grs *gameroomslice) JoinRoom(stream *rpsSvcGameServer) *gameroom {
	grs.Lock()
	defer grs.Unlock()
	return grs.joinRoom(stream)
}

// joinRoom : returns a gameroom from stack or creates one if there aren't any.
func (grs *gameroomslice) joinRoom(stream *rpsSvcGameServer) *gameroom {
	var g *gameroom
	grsLen := len(grs.slRooms)
	if grsLen < 1 {
		// no rooms in slice
		// create and add to slice
		g = &gameroom{Player1: stream, IsFull: make(chan bool)}
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
