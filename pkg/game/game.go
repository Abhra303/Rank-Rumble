package game

import (
	"log"
	"net/http"
	"sync"

	"github.com/efficientgo/core/errors"
	"github.com/gorilla/websocket"
)

type Card int

type PlayerEvent struct {
}

type playerResponse struct {
	Err    string `json:"error"`
	Winner string `json:"winner"`
	Status int    `json:"status"`
}

type Player struct {
	Conn       *websocket.Conn
	totalCards int
	emitEvent  chan PlayerEvent
}

type EndInfo struct {
	Err    error
	Winner int
}

type Game struct {
	drawPile    []Card
	discardPile []Card
	players     []Player
	playerTurn  int

	winnerSignal chan int
	EndSignal    chan EndInfo
	eventEmitter chan playerResponse
	Err          error
	wg           *sync.WaitGroup
}

func (g *Game) Start() {
	g.wg.Add(1)
	go g.broadcastPlayEventListener()
	go g.listenPlayEvent()
	select {
	case end := <-g.EndSignal:
		var event playerResponse
		if end.Err != nil {
			event.Err = end.Err.Error()
			g.broadcastPlayEvent(event)
			return
		}
		event.Winner = g.players[end.Winner].Conn.LocalAddr().String()
		g.broadcastPlayEvent(event)
		close(g.eventEmitter)
		close(g.winnerSignal)
		return
	default:
	}
	g.wg.Wait()
}

func (g *Game) GetPlayer(index int) (Player, error) {
	if index >= len(g.players) {
		return Player{}, errors.New("given index exceeds the player index")
	}
	return g.players[index], nil
}

func (g *Game) getWinner() {
	c := <-g.winnerSignal
	var err error
	if c >= len(g.players) {
		err = errors.New("winner index exceeds the player index")
	}
	g.EndSignal <- EndInfo{Err: err, Winner: c}
	close(g.EndSignal)
}

func (g *Game) broadcastPlayEventListener() {
	for event := range g.eventEmitter {
		g.broadcastPlayEvent(event)
	}
}

func (g *Game) broadcastPlayEvent(event playerResponse) {
	for _, player := range g.players {
		go (func(c *websocket.Conn) {
			err := c.WriteJSON(&event)
			if err != nil {
				log.Println(err)
				return
			}
		})(player.Conn)
	}
}

func (g *Game) listenPlayEvent() {
	defer g.wg.Done()
	for {
		var m PlayerEvent
		err := g.players[g.playerTurn].Conn.ReadJSON(&m)
		if err != nil {
			log.Println(err)
			g.eventEmitter <- playerResponse{Err: err.Error(), Status: http.StatusBadRequest}
			return
		}
		// code for game condition
		// TODO: implement the code
		g.EndSignal <- EndInfo{Err: nil, Winner: 0}
		close(g.EndSignal)
		break
	}
}

func NewGame(players []Player) (*Game, error) {
	if len(players) < 2 {
		return nil, errors.New("minimum 2 players needed to start a game")
	}
	return &Game{
		players:      players,
		drawPile:     make([]Card, 0, 12),
		discardPile:  make([]Card, 0, 5),
		winnerSignal: make(chan int),
		EndSignal:    make(chan EndInfo),
		wg:           new(sync.WaitGroup),
	}, nil
}
