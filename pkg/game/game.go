package game

import (
	"github.com/efficientgo/core/errors"
	"github.com/gorilla/websocket"
)

type Card int

type PlayerEvent struct {
}

type Player struct {
	Conn       *websocket.Conn
	totalCards int
	emitEvent  chan PlayerEvent
}

type EndInfo struct {
	err    error
	winner int
}

type game struct {
	drawPile    []Card
	discardPile []Card
	players     []Player

	winnerSignal chan int
	EndSignal    chan EndInfo
	Err          error
}

type Game interface {
	Start()
	GetPlayer(int) (Player, error)
	getWinner()
}

func (g *game) Start() {}

func (g *game) GetPlayer(index int) (Player, error) {
	if index >= len(g.players) {
		return Player{}, errors.New("given index exceeds the player index")
	}
	return g.players[index], nil
}

func (g *game) getWinner() {
	c := <-g.winnerSignal
	var err error
	if c >= len(g.players) {
		err = errors.New("winner index exceeds the player index")
	}
	g.EndSignal <- EndInfo{err: err, winner: c}
}

func NewGame(players []Player) (Game, error) {
	if len(players) < 2 {
		return nil, errors.New("minimum 2 players needed to start a game")
	}
	return &game{
		players:      players,
		drawPile:     make([]Card, 0, 12),
		discardPile:  make([]Card, 0, 5),
		winnerSignal: make(chan int),
		EndSignal:    make(chan EndInfo),
	}, nil
}
