package room

import (
	"sync"

	"github.com/Abhra303/Rank-Rumble/pkg/game"
	"github.com/efficientgo/core/errors"
	"github.com/gorilla/websocket"
)

const DefaultRoomLimit = 4

type room struct {
	conns          []*websocket.Conn
	maxPlayerLimit int
	currentPlayers int
	mut            sync.Mutex
}

type Room interface {
	JoinRoom(*websocket.Conn) (bool, error)
	MaxPlayerLimit() int
	LeaveRoom(*websocket.Conn) (bool, error)
	CurrentPlayersSize() int
	StartGame() error
	IsFull() bool
}

func (r *room) JoinRoom(conn *websocket.Conn) (bool, error) {
	if conn == nil {
		return false, errors.New("nil connection can't join a room")
	}
	r.mut.Lock()
	defer r.mut.Unlock()
	if !r.IsFull() {
		r.conns = append(r.conns, conn)
		return true, nil
	}
	return false, nil
}

func (r *room) LeaveRoom(conn *websocket.Conn) (bool, error) {
	if conn == nil {
		return false, errors.New("nil connection can't leave a room")
	}

	r.mut.Lock()
	defer r.mut.Unlock()

	if len(r.conns) == 1 {
		r.conns = nil
		r.currentPlayers = 0
		return true, nil
	}

	for i := range r.conns {
		if r.conns[i] == conn {
			r.conns = append(r.conns[:i], r.conns[i+1:]...)
			return false, nil
		}
	}
	return false, errors.New("the given connection doesn't belong to the room")
}

func (r *room) MaxPlayerLimit() int {
	return r.maxPlayerLimit
}

func (r *room) CurrentPlayersSize() int {
	return r.currentPlayers
}

func (r *room) IsFull() bool {
	if r.MaxPlayerLimit() == 0 {
		return false
	}
	return r.MaxPlayerLimit() <= r.CurrentPlayersSize()
}

func NewRoom(conn *websocket.Conn, maxPlayer int) Room {
	rm := new(room)
	rm.maxPlayerLimit = maxPlayer
	rm.conns = make([]*websocket.Conn, 0, DefaultRoomLimit)
	rm.conns = append(rm.conns, conn)
	rm.currentPlayers = 1
	return rm
}

func (r *room) StartGame() error {
	if r.conns == nil {
		return errors.New("room doesn't have any players")
	}
	players := make([]game.Player, len(r.conns))
	for i, conn := range r.conns {
		players[i] = game.Player{Conn: conn}
	}
	g, err := game.NewGame(players)
	if err != nil {
		return err
	}
	g.Start()
	return nil
}
