package room_test

import (
	"testing"

	"github.com/Abhra303/Rank-Rumble/pkg/room"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func test_prepare_rooms(serial bool, length int) []room.Room {
	rooms := make([]room.Room, 0, 4)
	for i := 0; i < length; i++ {
		var rm room.Room
		if !serial {
			rm, _ = room.NewRoom(&websocket.Conn{}, 4)
		} else {
			rm, _ = room.NewRoom(&websocket.Conn{}, i+1)
		}
		rooms = append(rooms, rm)
	}
	return rooms
}

func TestRoomPool(t *testing.T) {
	var err error
	rooms := test_prepare_rooms(false, 4)

	rp := room.NewRoomPool()
	for _, rm := range rooms {
		err = rp.AddRoom(rm)
		assert.NoError(t, err)
	}
	err = rp.AddRoom(rooms[0])
	assert.EqualError(t, err, "given room is already in the room pool")
	err = rp.AddRoom(nil)
	assert.EqualError(t, err, "the given room is not valid")

	for _, rm := range rooms {
		err = rp.RemoveRoom(rm)
		assert.NoError(t, err)
	}
	rm, _ := room.NewRoom(&websocket.Conn{}, 4)
	err = rp.RemoveRoom(rm)
	assert.EqualError(t, err, "the given room is not found in this room pool")
	err = rp.RemoveRoom(nil)
	assert.EqualError(t, err, "the given room is not valid")
}

func TestMatchRoom(t *testing.T) {
	var err error
	rooms := test_prepare_rooms(true, 5)

	rp := room.NewRoomPool()
	for _, rm := range rooms {
		err = rp.AddRoom(rm)
		assert.NoError(t, err)
	}

	rm, err := rp.MatchRoom(&websocket.Conn{}, 2)
	assert.NoError(t, err)
	assert.NotNil(t, rm)
	assert.Equal(t, 2, rm.MaxPlayerLimit())
	assert.True(t, rm.IsFull())

	rm, err = rp.MatchRoom(&websocket.Conn{}, 2)
	assert.NoError(t, err)
	assert.Nil(t, rm)

	rm, err = rp.MatchRoom(&websocket.Conn{}, 3)
	assert.NoError(t, err)
	assert.NotNil(t, rm)
	assert.Equal(t, 2, rm.CurrentPlayersSize())
	assert.Equal(t, 3, rm.MaxPlayerLimit())

	rm, err = rp.MatchRoom(nil, 5)
	assert.EqualError(t, err, "connection can't be nil")
	assert.Nil(t, rm)
}
