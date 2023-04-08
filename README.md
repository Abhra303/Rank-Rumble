# Rank Rumble: The multiplayer card game

Rank Rumble is a multiplayer game backend project completely written in Go. Due to Go's smooth concurrency ecosystem, the system is light-weight and fast. The server run on localhost at port `8080`.

## Rules of the Game:

* Each player starts with a hand of 5 cards.
* The game starts with a deck of 52 cards ( a standard deck of playing cards).
* Players take turns playing cards from their hand, following a set of rules that define what cards can be played when.
* A player can only play a card if it matches either the suit or the rank of the top card on the discard pile.
* If a player cannot play a card, they must draw a card from the draw pile. If the draw pile is empty, the game ends in a draw and no player is declared a winner.
* The game ends when one player runs out of cardswho is declared the winner.

## How to run

To run the project, clone the repo using `git clone github.com/Abhra303/Rank-Rumble.git`. Go to Rank-Rumble directory and run `go run cmd/*.go`.

```
$ git clone github.com/Abhra303/Rank-Rumble.git

$ cd Rank-Rumble

$ go run cmd/*.go
```

If you have `make` installed in your machine, you can simply run `make run` command to run the server. There is also a `make test` command which can be used to run the unit tests.

```
$ make test

go test ./...
?       github.com/Abhra303/Rank-Rumble/cmd     [no test files]
ok      github.com/Abhra303/Rank-Rumble/pkg/game        0.999s
ok      github.com/Abhra303/Rank-Rumble/pkg/room        (cached)
```

## Architecture

The architecture is simple yet powerful. Players will first use http request to connect to the server. The connection will then upgrade into a socket connection after a dual handshake.

Players can create a new "Room" or join an existing Room. Player who created a particular room can specify the maximum limit of players who can join the room. Default limit of players in a room is 4. You can extend or decrease the limit. Note that atleast 2 players needed in a room to start the game.

Structurally, there are three steps to perform the operations. These are - (i) Client socket connection setup (ii) Room creation and joining (iii) Game engine initialization and perform game operations.

I am using JSON for socket communications between client and server. It would be better if I could use `protobuf` here (May be in the v2 version). But as the number of players in a card game is smaller, we can perfectly go with the JSON.

There are two types of requests/responses json structure. The first type of requests/response structure is used in the "client" step where the communications are used to create/join rooms and starting the game.

Below is the structure of first type of requests/response:

**Request body structure**

Note that you can omit any fields depending on what you want to do.

```json
{
	"maxPlayerLimit": 4,
	"createRoom": false,
	"joinRoomIfMatch": false,
	"startGame": false,
	"endGame": false
}
```

The above structure will create a room with player limit 4 and include the client in the room.

**Response body structure**

```json
{
	"error": "",
	"status": 200
}
```

The second type of request/response structure is used in "game" step where the communication are used to control various events of the game (e.g. getting played card information by a player).

**Request body structure**

```json
{
	"card": {
		"value": 6,
		"suit": 2
	},
	"skip": false
}
```

Note that "skip" is optional here.

**Response body structure**

```json
{
	"error": "",
	"status": 200,
	"draw": true,
	"winner": ""
}
```

The server is fast and multiple concurrent routines are used to optimize and improve the performance. I used [stretchr/testify](github.com/stretchr/testify) for unit testing and [gorilla/websocket](github.com/gorilla/websocket) for socket handling.
