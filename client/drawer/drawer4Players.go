package drawer

import "fmt"

// init some variables for 4-players game
func initFourPlayers(index int32) {
	initFourPlayersTemplate()
	Draw = drawFourPlayers
	SetDecks = setDecksFourPlayers
	indexString = indexStringFourPlayers

	arrowStates[0] = "↓"
	arrowStates[1] = "←"
	arrowStates[2] = "↑"
	arrowStates[3] = "→"

	indexesStringNumber = 26

	center = index
	left = (index + 1) % 4
	opposite = (index + 2) % 4
	right = (index + 3) % 4
}

// SetDecks for 4-players game
func setDecksFourPlayers(cards []int32) {
	for i := 0; i < len(cards); i++ {
		deck[center][i] = fmt.Sprintf("%s%s", suits[cards[i]%4], ranks[cards[i]/4])
		deck[left][i] = closedCard
		deck[opposite][i] = closedCard
		deck[right][i] = closedCard
	}
}

// indexString for 4-players game
func indexStringFourPlayers(n int) string {
	str := emptyString(17)
	for i := 1; i <= n; i++ {
		str += fmt.Sprintf("%3d", i)
	}
	str += emptyString(17 + (9-n)*3)
	return str
}

// inits template for 4-players game
func initFourPlayersTemplate() {
	template = make([]string, 28)
	template[0] = emptyString(17) + cardString(9) + emptyString(17)

	template[1] = emptyString(61)

	for i := 2; i <= 4; i++ {
		template[i] = emptyString(5) + infoString() + emptyString(39)
	}

	template[5] = emptyString(61)

	template[6] = emptyString(30) + "%1s" + emptyString(30)

	template[7] = emptyString(61)

	for i := 8; i <= 9; i++ {
		template[i] = emptyString(3) + infoString() + emptyString(21) + infoString() + emptyString(3)
	}

	template[10] = emptyString(3) + infoString() + emptyString(9) + "___" + emptyString(9) + infoString() + emptyString(3)

	template[11] = emptyString(2) + cardString(1) + emptyString(23) + "|" + cardString(1) + "|" + emptyString(23) + cardString(1) + emptyString(2)

	template[12] = emptyString(2) + cardString(1) + emptyString(24) + "‾‾‾" + emptyString(24) + cardString(1) + emptyString(2)

	template[13] = emptyString(2) + cardString(1) + emptyString(51) + cardString(1) + emptyString(2)

	template[14] = emptyString(2) + cardString(1) + emptyString(10) + "___" + emptyString(25) + "___" + emptyString(10) + cardString(1) + emptyString(2)

	template[15] = emptyString(2) + cardString(1) + emptyString(9) + "|" + cardString(1) + "|" + emptyString(11) + "%1s" + emptyString(11) + "|" +
		cardString(1) + "|" + emptyString(9) + cardString(1) + emptyString(2)

	template[16] = emptyString(2) + cardString(1) + emptyString(10) + "‾‾‾" + emptyString(25) + "‾‾‾" + emptyString(10) + cardString(1) + emptyString(2)

	template[17] = emptyString(2) + cardString(1) + emptyString(51) + cardString(1) + emptyString(2)

	template[18] = emptyString(2) + cardString(1) + emptyString(24) + "___" + emptyString(24) + cardString(1) + emptyString(2)

	template[19] = emptyString(2) + cardString(1) + emptyString(23) + "|" + cardString(1) + "|" + emptyString(23) + cardString(1) + emptyString(2)

	template[20] = emptyString(29) + "‾‾‾" + emptyString(29)

	template[21] = emptyString(61)

	for i := 22; i <= 24; i++ {
		template[i] = emptyString(39) + infoString() + emptyString(5)
	}

	template[25] = emptyString(17) + cardString(9) + emptyString(17)

	template[26] = emptyString(61)
	template[27] = emptyString(61)

	for i := 0; i < 28; i++ {
		template[i] += "\n"
	}
}

// Draw for 4-players game
func drawFourPlayers() {
	drawSpace()

	fmt.Printf(template[0], deck[opposite][0], deck[opposite][1], deck[opposite][2], deck[opposite][3], deck[opposite][4],
		deck[opposite][5], deck[opposite][6], deck[opposite][7], deck[opposite][8])
	fmt.Printf(template[1])
	fmt.Printf(template[2], pointsString, score[opposite])
	fmt.Printf(template[3], curBetString, bets[opposite])
	fmt.Printf(template[4], curBoardsString, wins[opposite])
	fmt.Printf(template[5])
	fmt.Printf(template[6], trump)
	fmt.Printf(template[7])
	fmt.Printf(template[8], pointsString, score[left], pointsString, score[right])
	fmt.Printf(template[9], curBetString, bets[left], curBetString, bets[right])
	fmt.Printf(template[10], curBoardsString, wins[left], curBoardsString, wins[right])
	fmt.Printf(template[11], deck[left][8], cards[opposite], deck[right][8])
	fmt.Printf(template[12], deck[left][7], deck[right][7])
	fmt.Printf(template[13], deck[left][6], deck[right][6])
	fmt.Printf(template[14], deck[left][5], deck[right][5])
	fmt.Printf(template[15], deck[left][4], cards[left], arrow, cards[right], deck[right][4])
	fmt.Printf(template[16], deck[left][3], deck[right][3])
	fmt.Printf(template[17], deck[left][2], deck[right][2])
	fmt.Printf(template[18], deck[left][1], deck[right][1])
	fmt.Printf(template[19], deck[left][0], cards[center], deck[right][0])
	fmt.Printf(template[20])
	fmt.Printf(template[21])
	fmt.Printf(template[22], pointsString, score[center])
	fmt.Printf(template[23], curBetString, bets[center])
	fmt.Printf(template[24], curBoardsString, wins[center])
	fmt.Printf(template[25], deck[center][0], deck[center][1], deck[center][2], deck[center][3], deck[center][4],
		deck[center][5], deck[center][6], deck[center][7], deck[center][8])
	fmt.Printf(template[26])
	fmt.Printf(template[27])
}
