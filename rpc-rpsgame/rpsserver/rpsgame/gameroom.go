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

type player struct {
	ID     string
	Stream *rpsSvcGameServer
}

type gameroom struct {
	RoomID  string
	Player1 *player
	Player2 *player
}

// GetRoom : thread safe version for getRoom
func (grs *gameroomslice) GetRoom() *gameroom {
	grs.Lock()
	defer grs.Unlock()
	return grs.getRoom()
}

// getRoom : returns a gameroom from stack or creates one if there aren't any.
func (grs *gameroomslice) getRoom() *gameroom {
	var g *gameroom
	grsLen := len(grs.slRooms)
	if grsLen < 1 {
		// no rooms in slice
		// create and add to slice
		g = &gameroom{}
		grs.slRooms = append(grs.slRooms, g)
	} else {
		// get first in
		// check if room has 2 both players
		// remove from slice
		g, grs.slRooms = grs.slRooms[0], grs.slRooms[1:]
	}
	return g
}
