package room

import (
	"sync"

	"github.com/efficientgo/core/errors"
	"github.com/gorilla/websocket"
)

// TODO: use map instead of array where the key would be
// the room id and value would be Room
type roomPool struct {
	rooms []Room
	mut   sync.Mutex
}

type RoomPool interface {
	AddRoom(Room) error
	RemoveRoom(Room) error
	MatchRoom(*websocket.Conn, int) (Room, error)
}

func (rp *roomPool) AddRoom(rm Room) error {
	if !IsRoomValid(rm) {
		return errors.New("the given room is not valid")
	}
	rp.mut.Lock()
	defer rp.mut.Unlock()

	// see if the room is already present
	for i := range rp.rooms {
		if rp.rooms[i] == rm {
			return errors.New("given room is already in the room pool")
		}
	}
	rp.rooms = append(rp.rooms, rm)
	return nil
}

func (rp *roomPool) RemoveRoom(rm Room) error {
	if !IsRoomValid(rm) {
		return errors.New("the given room is not valid")
	}

	rp.mut.Lock()
	defer rp.mut.Unlock()

	for i := range rp.rooms {
		if rp.rooms[i] == rm {
			rp.rooms = append(rp.rooms[:i], rp.rooms[i+1:]...)
			return nil
		}
	}
	return errors.New("the given room is not found in this room pool")
}

// TODO: add a timer so that it can wait for new room creation
func (rp *roomPool) MatchRoom(conn *websocket.Conn, maxPlayerLength int) (Room, error) {
	if conn == nil {
		return nil, errors.New("connection can't be nil")
	}

	rp.mut.Lock()
	defer rp.mut.Unlock()

	for _, room := range rp.rooms {
		if room.CurrentPlayersSize() != 0 && room.MaxPlayerLimit() == maxPlayerLength && !room.IsFull() {
			_, err := room.JoinRoom(conn)
			if err != nil {
				return nil, err
			}
			return room, nil
		}
	}
	return nil, nil
}

func NewRoomPool() RoomPool {
	return &roomPool{
		rooms: make([]Room, 0, 4),
	}
}
