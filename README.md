# Rank Rumble: The multiplayer card game

Rank Rumble is a multiplayer game backend project completely written in Go. Due to Go's smooth concurrency ecosystem, the system is light-weight and fast.

## Rules of the Game:

* Each player starts with a hand of 5 cards.
* The game starts with a deck of 52 cards ( a standard deck of playing cards).
* Players take turns playing cards from their hand, following a set of rules that define what cards can be played when.
* A player can only play a card if it matches either the suit or the rank of the top card on the discard pile.
* If a player cannot play a card, they must draw a card from the draw pile. If the draw pile is empty, the game ends in a draw and no player is declared a winner.
* The game ends when one player runs out of cardswho is declared the winner.

## How to run

To run the project, clone the repo using `git clone github.com/Abhra303/Rank-Rumble.git`. Go to Rank-Rumble directory and run `go run cmd/*.go`.

```bash
$ git clone github.com/Abhra303/Rank-Rumble.git

$ cd Rank-Rumble

$ go run cmd/*.go
```
