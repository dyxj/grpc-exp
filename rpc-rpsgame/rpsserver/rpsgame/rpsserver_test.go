package rpsgame

import (
	"fmt"
	"log"
	"testing"
)

func TestNewRpsSvc(t *testing.T) {
	rspsvc, err := NewRpsSvc()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rspsvc.rooms.slRooms)
}
