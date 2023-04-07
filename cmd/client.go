package main

import (
	"log"

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
}

type client struct {
	conn     *websocket.Conn
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
			return
		}
		if m.CreateRoom {
			if m.MaxPlayerLimit == 0 {
				m.MaxPlayerLimit = room.DefaultRoomLimit
			}
			// code for room creation
		} else if m.JoinRoomIfMatch {
			// for each elem in room see the condition
			// if satisy, join into the room
		} else if m.StartGame {
			// code for starting the game
		} else if m.EndGame {
			// code for ending game
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
	return &client{conn: conn, writeMsg: make(chan responseMessage, 1)}
}
