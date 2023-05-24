// package drawer provides a console client for hussar.
// Functions are:
// func Init       (int32, int32)
// func UpdateIndex(int32)
// func SetBet     (int32, int32)
// func SetDecks   ([]int32)
// func SetCard    (int32, int32)
// func SetTrump   (int32)
// func SetArrow   (int32)
// func FlushCard  (int32, int32)
// func FlushBoard ()
// func FlushRound ()
// func AddWin     (int32)
// func AddRes     (int32, int32)
// func ReplaceCard(int32, int32)
// func Draw       ()
package drawer

import (
	"fmt"
)

// Const strings to be written on the screen
const (
	pointsString    = "Points"
	curBetString    = "Boards bet"
	curBoardsString = "Boards won"
)

// Const strings to represent cards
const (
	closedCard = "🂠"
	empty      = " "
)

// Variables, which state is displayed on the board
var (
	template                      []string   // template for drawer
	playersN                      int        // number of players
	left, right, center, opposite int32      // indexex of players. Center is always client's index
	deck                          [][]string // players' decks
	cards                         []string   // cards, which are on the board
	score, bets, wins             []int32    // numbers, which are shown during the game
	trump                         string     // trump of current board
	arrow                         string     // arrow points whose turn to bet/move now
	arrowStates                   []string   // possible states of arrow
	indexesStringNumber           int        // number of string with indexes, which helps to choose card
)

// Init inits a template with total as a number of players and index as a client index
// Must be called before any other drawer function
func Init(total, index int32) {
	initSlices(total)
	playersN = int(total)

	if total == 3 {
		initThreePlayers(index)
	} else if total == 2 {
		initTwoPlayers(index)
	} else if total == 4 {
		initFourPlayers(index)
	}
}

// initSlices inits slices to store some info with total as a number of players
func initSlices(total int32) {
	arrowStates = make([]string, total)
	cards = make([]string, total)
	score = make([]int32, total)
	bets = make([]int32, total)
	wins = make([]int32, total)

	deck = make([][]string, total)
	for i := int32(0); i < total; i++ {
		deck[i] = make([]string, 36/total)
	}
}

// UpdateIndex updates index string,
// which is displayed to ease choice which card to use
func UpdateIndex(n int32) {
	template[indexesStringNumber] = indexString(int(n))
}

// SetBet sets one's bet
func SetBet(index, bet int32) {
	bets[index] = bet
}

// SetDecks sets cards in the beginning of a round
// Has different implementation for different number of players.
var SetDecks func([]int32)

// SetCard sets card on the board when some player uses it
func SetCard(index, c int32) {
	cards[index] = fmt.Sprintf("%s%s", suits[c%4], ranks[c/4])
}

// SetTrump sets trump in the beginning of a round
func SetTrump(tr int32) {
	trump = suits[tr]
}

// SetArrow updates the arrow so that it points to player, who moves now
func SetArrow(i int32) {
	arrow = arrowStates[i]
}

// FlushCard flushes a card when player uses it
func FlushCard(index, card int32) {
	deck[index][card] = empty
}

// FlushBoard flushes the information of the past board
func FlushBoard() {
	for i := 0; i < playersN; i++ {
		cards[i] = empty
	}
}

// FlushRound flushes the information of the past round
func FlushRound() {
	for i := 0; i < playersN; i++ {
		bets[i] = 0
		wins[i] = 0
	}
}

// AddWin adds a point to player who won
func AddWin(index int32) {
	wins[index]++
}

// AddRes updates total score of players
func AddRes(index, res int32) {
	score[index] += res
}

// ReplaceCard replaces player's card with index "to" with card with index "with"
func ReplaceCard(to, from int32) {
	deck[center][to] = deck[center][from]
}

// Draw displays the board.
// Has different implementation for different number of players.
var Draw func()

// drawSpace draws empty lines to visually divide different Draw()
func drawSpace() {
	for i := 0; i < 50; i++ {
		fmt.Println()
	}
}

// cardString returns a place which will be taken by card for the template
// n stands for number of cards
func cardString(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += "%3s"
	}
	return s
}

// emptyString returns a string of n spaces for the template
func emptyString(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += " "
	}
	return s
}

// infoString returns a string to be filled with points/bet info for template
func infoString() string {
	return "%12s%5d"
}

// indexString returns a string for template with numbers 1..n written
var indexString func(int) string
