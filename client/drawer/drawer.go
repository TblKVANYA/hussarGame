// package drawer provides a console client for hussar.
// Functions are:
// func Draw()
// func UpdateIndex(int32)
// func SetSides   (int32)
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
	left, right, center int32         // indexex of players. Center = own index
	deck                [3][12]string // players' decks
	cards               [3]string     // cards, which are on the board
	score, bets, wins   [3]int32      // numbers, which are shown during the game
	trump               string        // trump of current board
	arrow               string        // arrow points whose turn to bet/move now
)

var arrowStates = [3]string{"↓", "←", "→"}

var template []string

// init inits a template
func init() {
	template = make([]string, 24)
	for i := 0; i < 3; i++ {
		template[i] = " " + infoString() + emptyString(9) + infoString() + " "
	}

	template[3] = emptyString(22) + "%1s" + emptyString(22)

	for i := 4; i <= 8; i++ {
		template[i] = " %1s" + emptyString(41) + "%1s "
	}
	template[9] = " %1s" + emptyString(20) + "%1s" + emptyString(20) + "%1s "

	template[10] = " %1s" + emptyString(10) + "___" + emptyString(15) + "___" + emptyString(10) + "%1s "

	template[11] = " %1s" + emptyString(9) + "|" + cardString(1) + "|" + emptyString(13) + "|" + cardString(1) + "|" + emptyString(9) + "%1s "

	template[12] = " %1s" + emptyString(10) + "‾‾‾" + emptyString(15) + "‾‾‾" + emptyString(10) + "%1s "

	for i := 13; i <= 14; i++ {
		template[i] = " %1s" + emptyString(41) + "%1s "
	}

	template[15] = " %1s" + emptyString(19) + "___" + emptyString(19) + "%1s "
	template[16] = emptyString(20) + "|" + cardString(1) + "|" + emptyString(20)
	template[17] = emptyString(21) + "‾‾‾" + emptyString(3) + infoString() + " "

	for i := 18; i <= 19; i++ {
		template[i] = emptyString(27) + infoString() + " "
	}

	template[20] = emptyString(45)
	template[21] = emptyString(4) + cardString(12) + emptyString(5)
	template[22] = emptyString(45)
	template[23] = emptyString(45)

	for i := 0; i < 24; i++ {
		template[i] += "\n"
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

// indexString returns a string with numbers 1..n written
func indexString(n int) (str string) {
	str = emptyString(4)
	for i := 1; i <= n; i++ {
		str = fmt.Sprintf("%s%3d", str, i)
	}
	str += emptyString(5 + (12-n)*3)
	return

}

// Draw displays the board
func Draw() {
	for i := 0; i < 50; i++ {
		fmt.Println()
	}

	fmt.Printf(template[0], pointsString, score[left], pointsString, score[right])
	fmt.Printf(template[1], curBetString, bets[left], curBetString, bets[right])
	fmt.Printf(template[2], curBoardsString, wins[left], curBoardsString, wins[right])
	fmt.Printf(template[3], trump)
	fmt.Printf(template[4], deck[left][11], deck[right][11])
	fmt.Printf(template[5], deck[left][10], deck[right][10])
	fmt.Printf(template[6], deck[left][9], deck[right][9])
	fmt.Printf(template[7], deck[left][8], deck[right][8])
	fmt.Printf(template[8], deck[left][7], deck[right][7])
	fmt.Printf(template[9], deck[left][6], arrow, deck[right][6])
	fmt.Printf(template[10], deck[left][5], deck[right][5])
	fmt.Printf(template[11], deck[left][4], cards[left], cards[right], deck[right][4])
	fmt.Printf(template[12], deck[left][3], deck[right][3])
	fmt.Printf(template[13], deck[left][2], deck[right][2])
	fmt.Printf(template[14], deck[left][1], deck[right][1])
	fmt.Printf(template[15], deck[left][0], deck[right][0])
	fmt.Printf(template[16], cards[center])
	fmt.Printf(template[17], pointsString, score[center])
	fmt.Printf(template[18], curBetString, bets[center])
	fmt.Printf(template[19], curBoardsString, wins[center])
	fmt.Printf(template[20])
	fmt.Printf(template[21], deck[center][0], deck[center][1], deck[center][2], deck[center][3], deck[center][4],
		deck[center][5], deck[center][6], deck[center][7], deck[center][8], deck[center][9], deck[center][10], deck[center][11])
	fmt.Printf(template[22])
	fmt.Printf(template[23])
}

// UpdateIndex updates index string,
// which is displayed to ease choice which card to use
func UpdateIndex(n int32) {
	template[22] = indexString(int(n))
}

// SetSides sets indexex of player and who sits to the left and right of him. Should be used when player's index is known
func SetSides(index int32) {
	center = index
	left = (index + 1) % 3
	right = (index + 2) % 3
}

// SetBet sets one's bet
func SetBet(index, bet int32) {
	bets[index] = bet
}

// SetDecks sets cards in the beginning of a round
func SetDecks(cards []int32) {
	for i := 0; i < len(cards); i++ {
		deck[center][i] = fmt.Sprintf("%s%s", suits[cards[i]%4], ranks[cards[i]/4])
		deck[left][i] = closedCard
		deck[right][i] = closedCard
	}
}

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
	for i := 0; i < 3; i++ {
		cards[i] = empty
	}
}

// FlushRound flushes the information of the past round
func FlushRound() {
	for i := 0; i < 3; i++ {
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
