package rpsgame

import (
	"fmt"
	"testing"
)

func TestGetRoom(t *testing.T) {
	grs := gameroomslice{}
	fmt.Println(grs.slRooms)
	gr := grs.GetRoom()
	fmt.Println(grs.slRooms)
	fmt.Println(gr)
}
