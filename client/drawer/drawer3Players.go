package drawer

import "fmt"

// init some variables for 3-players game
func initThreePlayers(index int32) {
	initThreePlayersTemplate()
	Draw = drawThreePlayers
	SetDecks = setDecksThreePlayers
	indexString = indexStringThreePlayers

	arrowStates[0] = "↓"
	arrowStates[1] = "←"
	arrowStates[2] = "→"

	indexesStringNumber = 22

	center = index
	left = (index + 1) % 3
	right = (index + 2) % 3
}

// SetDecks for 3-players game
func setDecksThreePlayers(cards []int32) {
	for i := 0; i < len(cards); i++ {
		deck[center][i] = fmt.Sprintf("%s%s", suits[cards[i]%4], ranks[cards[i]/4])
		deck[left][i] = closedCard
		deck[right][i] = closedCard
	}
}

// indexString for 3-players game
func indexStringThreePlayers(n int) string {
	str := emptyString(4)
	for i := 1; i <= n; i++ {
		str += fmt.Sprintf("%3d", i)
	}
	str += emptyString(5 + (12-n)*3)
	return str
}

// inits template for 3-players game
func initThreePlayersTemplate() {
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

// Draw for 3-players game
func drawThreePlayers() {
	drawSpace()

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
