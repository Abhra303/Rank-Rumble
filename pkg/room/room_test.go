package room_test

import (
	"testing"

	"github.com/Abhra303/Rank-Rumble/pkg/room"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestRoomPlayerLimit(t *testing.T) {
	tcs := []struct {
		name           string
		maxPlayerLimit int
		shouldErr      bool
		err            string
	}{
		{
			name:           "maxPlayerLimit is negative",
			maxPlayerLimit: -1,
			shouldErr:      true,
			err:            "max player limit can't be negative or zero",
		},
		{
			name:           "maxPlayerLimit zero returns error",
			maxPlayerLimit: 0,
			shouldErr:      true,
			err:            "max player limit can't be negative or zero",
		},
		{
			name:           "maxPlayerLimit is one",
			maxPlayerLimit: 1,
		},
		{
			name:           "maxPlayerLimit is 50",
			maxPlayerLimit: 50,
		},
	}

	for _, tt := range tcs {
		t.Run(tt.name, func(t *testing.T) {
			var conn *websocket.Conn
			room, err := room.NewRoom(conn, tt.maxPlayerLimit)
			if tt.shouldErr {
				assert.EqualError(t, err, tt.err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, room)
			assert.Equal(t, tt.maxPlayerLimit, room.MaxPlayerLimit(), "they should be equal")
			assert.Equal(t, 1, room.CurrentPlayersSize(), "new room should have only one connection")
		})
	}
}

func TestRoomFunctionalities(t *testing.T) {
	conn1 := &websocket.Conn{}
	conn2 := &websocket.Conn{}
	conn3 := &websocket.Conn{}
	conn4 := &websocket.Conn{}
	rm, _ := room.NewRoom(conn1, 4)

	ok, err := rm.JoinRoom(conn2)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, 2, rm.CurrentPlayersSize())

	ok, err = rm.JoinRoom(conn3)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, 3, rm.CurrentPlayersSize())

	ok, err = rm.JoinRoom(conn4)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, 4, rm.CurrentPlayersSize())
	assert.True(t, rm.IsFull())

	ok, err = rm.JoinRoom(&websocket.Conn{})
	assert.NoError(t, err)
	assert.False(t, ok)

	empty, err := rm.LeaveRoom(conn1)
	assert.NoError(t, err)
	assert.False(t, empty)
	assert.Equal(t, 3, rm.CurrentPlayersSize())

	empty, err = rm.LeaveRoom(conn1)
	assert.Error(t, err)
	assert.False(t, empty)
	assert.Equal(t, 3, rm.CurrentPlayersSize())

	empty, err = rm.LeaveRoom(conn3)
	assert.NoError(t, err)
	assert.False(t, empty)
	assert.Equal(t, 2, rm.CurrentPlayersSize())

	empty, err = rm.LeaveRoom(conn4)
	assert.NoError(t, err)
	assert.False(t, empty)
	assert.Equal(t, 1, rm.CurrentPlayersSize())

	empty, err = rm.LeaveRoom(conn2)
	assert.NoError(t, err)
	assert.True(t, empty)
	assert.Equal(t, 0, rm.CurrentPlayersSize())

	assert.Equal(t, false, room.IsRoomValid(rm))
}
