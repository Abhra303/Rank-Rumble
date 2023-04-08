package main

import (
	"log"
	"net/http"

	"github.com/Abhra303/Rank-Rumble/pkg/room"
	"github.com/gorilla/websocket"
)

type requestMessage struct {
	MaxPlayerLimit  int  `json:"maxPlayerLimit,omitempty"`
	CreateRoom      bool `json:"createRoom,omitempty"`
	JoinRoomIfMatch bool `json:"joinRoomIfMatch,omitempty"`
	StartGame       bool `json:"startGame,omitempty"`
	EndGame         bool `json:"endGame,omitempty"`
}

type responseMessage struct {
	Err    string `json:"error"`
	Status int    `json:"status"`
}

type client struct {
	conn     *websocket.Conn
	room     room.Room
	writeMsg chan responseMessage
}

type Client interface {
	Listen()
	ClientInfo()
}

func (c *client) Listen() {
	go c.listenWrite()
	c.listenRead()
}

func (c *client) ClientInfo() {}

func (c *client) listenRead() {
	for {
		var m requestMessage
		err := c.conn.ReadJSON(&m)
		if err != nil {
			log.Println(err)
			c.writeMsg <- responseMessage{Err: err.Error(), Status: http.StatusBadRequest}
			return
		}
		if m.CreateRoom {
			// code for room creation
			if m.MaxPlayerLimit == 0 {
				m.MaxPlayerLimit = room.DefaultRoomLimit
			}
			rm := room.NewRoom(c.conn, m.MaxPlayerLimit)
			err = AvailableRooms.AddRoom(rm)
			if err != nil {
				log.Println(err)
				c.writeMsg <- responseMessage{Err: err.Error(), Status: http.StatusInternalServerError}
				return
			}
			c.room = rm
			c.writeMsg <- responseMessage{Err: "", Status: http.StatusCreated}
		} else if m.JoinRoomIfMatch {
			// for each elem in room see the condition
			// if satisy, join into the room
			if c.room != nil {
				c.writeMsg <- responseMessage{Err: "client already aligned with a room", Status: http.StatusBadRequest}
				return
			}
			if m.MaxPlayerLimit == 0 {
				m.MaxPlayerLimit = room.DefaultRoomLimit
			}

			rm, err := AvailableRooms.MatchRoom(c.conn, m.MaxPlayerLimit)
			if err != nil {
				log.Println(err)
				c.writeMsg <- responseMessage{Err: err.Error(), Status: http.StatusInternalServerError}
			}
			c.room = rm
			c.writeMsg <- responseMessage{Err: "", Status: http.StatusOK}
		} else if m.StartGame {
			// code for starting the game
			return
		} else if m.EndGame {
			// code for ending game
			return
		}
	}
}

func (c *client) listenWrite() {
	for {
		m := <-c.writeMsg
		err := c.conn.WriteJSON(&m)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func NewClient(conn *websocket.Conn) Client {
	return &client{conn: conn, writeMsg: make(chan responseMessage)}
}
