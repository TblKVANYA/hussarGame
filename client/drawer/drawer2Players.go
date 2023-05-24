package drawer

import "fmt"

// init some variables for 2-players game
func initTwoPlayers(index int32) {
	initTwoPlayersTemplate()
	Draw = drawTwoPlayers
	SetDecks = setDecksTwoPlayers
	indexString = indexStringTwoPlayers

	arrowStates[0] = "↓"
	arrowStates[1] = "↑"

	indexesStringNumber = 21
	center = index
	opposite = (index + 1) % 2

}

// SetDecks for 2-players game
func setDecksTwoPlayers(cards []int32) {
	for i := 0; i < len(cards); i++ {
		deck[center][i] = fmt.Sprintf("%s%s", suits[cards[i]%4], ranks[cards[i]/4])
		deck[opposite][i] = closedCard
	}
}

// indexString for 2-players game
func indexStringTwoPlayers(n int) string {
	str := emptyString(3)
	for i := 1; i <= n; i++ {
		str += fmt.Sprintf("%3d", i)
	}
	str += emptyString(4 + (18-n)*3)
	return str
}

// inits template for 2-players game
func initTwoPlayersTemplate() {
	template = make([]string, 23)
	template[0] = emptyString(3) + cardString(18) + emptyString(4)

	template[1] = emptyString(61)

	for i := 2; i < 5; i++ {
		template[i] = emptyString(3) + infoString() + emptyString(41)
	}

	template[5] = emptyString(29) + cardString(1) + emptyString(29)

	for i := 6; i < 10; i++ {
		template[i] = emptyString(61)
	}

	template[10] = emptyString(30) + "%1s" + emptyString(15) + "%1s" + emptyString(14)

	for i := 11; i < 15; i++ {
		template[i] = emptyString(61)
	}

	template[15] = emptyString(29) + cardString(1) + emptyString(29)

	for i := 16; i < 19; i++ {
		template[i] = emptyString(41) + infoString() + emptyString(3)
	}

	template[19] = emptyString(61)

	template[20] = emptyString(3) + cardString(18) + emptyString(4)

	template[21] = emptyString(61)
	template[22] = emptyString(61)

	for i := 0; i < 23; i++ {
		template[i] += "\n"
	}
}

// Draw for 2-players game
func drawTwoPlayers() {
	drawSpace()

	fmt.Printf(template[0], deck[opposite][0], deck[opposite][1], deck[opposite][2], deck[opposite][3], deck[opposite][4],
		deck[opposite][5], deck[opposite][6], deck[opposite][7], deck[opposite][8], deck[opposite][9], deck[opposite][10], deck[opposite][11],
		deck[opposite][12], deck[opposite][13], deck[opposite][14], deck[opposite][15], deck[opposite][16], deck[opposite][17])
	fmt.Printf(template[1])
	fmt.Printf(template[2], pointsString, score[opposite])
	fmt.Printf(template[3], curBetString, bets[opposite])
	fmt.Printf(template[4], curBoardsString, wins[opposite])
	fmt.Printf(template[5], cards[opposite])
	fmt.Printf(template[6])
	fmt.Printf(template[7])
	fmt.Printf(template[8])
	fmt.Printf(template[9])
	fmt.Printf(template[10], arrow, trump)
	fmt.Printf(template[11])
	fmt.Printf(template[12])
	fmt.Printf(template[13])
	fmt.Printf(template[14])
	fmt.Printf(template[15], cards[center])
	fmt.Printf(template[16], pointsString, score[center])
	fmt.Printf(template[17], curBetString, bets[center])
	fmt.Printf(template[18], curBoardsString, wins[center])
	fmt.Printf(template[19])
	fmt.Printf(template[20], deck[center][0], deck[center][1], deck[center][2], deck[center][3], deck[center][4],
		deck[center][5], deck[center][6], deck[center][7], deck[center][8], deck[center][9], deck[center][10], deck[center][11],
		deck[center][12], deck[center][13], deck[center][14], deck[center][15], deck[center][16], deck[center][17])
	fmt.Printf(template[21])
	fmt.Printf(template[22])
}
