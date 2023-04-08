package game

import (
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/efficientgo/core/errors"
	"github.com/gorilla/websocket"
)

type SuitType int

const (
	NoSuit SuitType = iota
	Club
	Heart
	Spade
	Diamond
)

type CardValue int

type Card struct {
	Value int
	Suit  SuitType
}

type PlayerEvent struct {
	Skip bool `json:"skip,omitempty"`
	Card Card `json:"card"`
}

type playerResponse struct {
	Err    string `json:"error"`
	Winner string `json:"winner"`
	Status int    `json:"status"`
}

type Player struct {
	Conn    *websocket.Conn
	GetData chan PlayerEvent
	Cards   []Card
}

type EndInfo struct {
	Err    error
	Winner int
	Draw   bool
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
}

func cardInPlayerCards(player Player, ccard Card) bool {
	for _, card := range player.Cards {
		if card.Value == ccard.Value {
			return true
		}
	}
	return false
}

func (g *Game) prepareCards() {
	deck := make([]Card, 52)
	suits := []SuitType{Club, Heart, Spade, Diamond}
	for i, suit := range suits {
		for j := 0; j < 13; j++ {
			deck[j+i*13] = Card{Value: j, Suit: suit}
		}
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]Card, 52)
	perm := r.Perm(52)
	for i, randIndex := range perm {
		ret[i] = deck[randIndex]
	}

	k := 0
	for i := 0; i < 5; i++ {
		for j, player := range g.players {
			k = i*4 + j
			player.Cards = append(player.Cards, ret[k])
		}
	}
	g.drawPile = append(g.drawPile, ret[k+1:]...)
}

func (g *Game) checkCardEligibility(cardPlayed PlayerEvent) bool {
	length := len(g.discardPile)
	if length == 0 {
		return true
	}
	dCard := g.discardPile[length-1]
	if cardPlayed.Card.Suit == dCard.Suit || cardPlayed.Card.Value >= dCard.Value {
		return true
	}
	return false
}

func (g *Game) Start() {
	go g.broadcastPlayEventListener()
	go g.listenPlayEvent()
	g.prepareCards()
	g.playerTurn = 0
	for {
		if g.playerTurn >= len(g.players) {
			g.playerTurn = 0
		}
		player := g.players[g.playerTurn]
		cardPlayed := <-player.GetData
		if !cardPlayed.Skip && !cardInPlayerCards(player, cardPlayed.Card) {
			g.eventEmitter <- playerResponse{Err: "the card is not present in player's stock"}
			continue
		}
		if cardPlayed.Skip || !g.checkCardEligibility(cardPlayed) {
			if len(g.drawPile) == 0 {
				g.EndSignal <- EndInfo{Draw: true}
			}
			card := g.drawPile[len(g.drawPile)-1]
			player.Cards = append(player.Cards, card)
			g.drawPile = g.drawPile[:len(g.drawPile)-1]
			continue
		}
		g.discardPile = append(g.discardPile, cardPlayed.Card)

		if len(player.Cards) == 0 {
			g.winnerSignal <- g.playerTurn
			close(g.winnerSignal)
			break
		}
		g.playerTurn++
	}
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
}

func (g *Game) GetPlayer(index int) (Player, error) {
	if index >= len(g.players) {
		return Player{}, errors.New("given index exceeds the player index")
	}
	return g.players[index], nil
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
	var wg *sync.WaitGroup
	wg.Add(len(g.players))
	for _, player := range g.players {
		go (func(player Player, wg *sync.WaitGroup) {
			defer wg.Done()
			for {
				var m PlayerEvent
				err := player.Conn.ReadJSON(&m)
				if err != nil {
					log.Println(err)
					g.eventEmitter <- playerResponse{Err: err.Error(), Status: http.StatusBadRequest}
					return
				}
				player.GetData <- m
			}
		})(player, wg)
	}
	wg.Wait()
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
	}, nil
}
