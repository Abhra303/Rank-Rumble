package main

import (
	"log"
	"net/http"

	"github.com/Abhra303/Rank-Rumble/pkg/room"
)

var AvailableRooms room.RoomPool

func main() {
	AvailableRooms = room.NewRoomPool()
	http.HandleFunc("/", CreateSocketHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
