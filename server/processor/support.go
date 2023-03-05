// Package processor implement the processor of the game, which encapsulates much of the logic of the program.
// The only callable function is Processor([3]datatypes.Tunnel)
package processor

import (
	"math/rand"

	"github.com/TblKVANYA/hussarGame/server/datatypes"
)

// holdRound holds the whole round with cardsN cards and returns updated table of players ans scores.
// Attacker stands for the player, who bets and moves on the first board of round before others.
// Fl can take values of 0 (usual round), 1 (dark round), 2 (untrumped round) and 3(golden round)
func holdRound(attacker, cardsN int32, table [3]int32, random *rand.Rand, tunnels [3]datatypes.Tunnel, fl int) [3]int32 {
	// Shuffle cards
	shuffled, trump := shuffleCards(random)

	// Is it untrumped round?
	if fl == 2 {
		trump = 4
	}

	// Is it dark round?
	df := int32(0)
	if fl == 1 {
		df = 1
	}
	// Send cards and round info
	sendCards(shuffled, cardsN, trump, attacker, tunnels, df)

	// Get bets
	bets := getBets(attacker, tunnels)

	// Hold boards
	wins := [3]int32{0, 0, 0}
	for n := int32(0); n < cardsN; n++ {
		attacker = holdBoard(attacker, trump, tunnels)
		wins[attacker]++
	}

	// Calculate results of round
	results := calculateResults(bets, wins)
	// Is it golden round?
	if fl == 3 {
		for i := 0; i < 3; i++ {
			results[i] *= 3
		}
	}
	// End round
	sendRes(results, tunnels)
	table = updateTable(table, results)

	return table
}

// Shuffle cards and randomly choose a trump.
func shuffleCards(random *rand.Rand) ([]datatypes.Card, int32) {
	cards := make([]datatypes.Card, datatypes.TotalCards)
	for i := int32(0); i < datatypes.TotalCards; i++ {
		cards[i] = datatypes.Card(i + 4)
	}

	random.Shuffle(int(datatypes.TotalCards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })

	trump := int32(cards[random.Intn(int(datatypes.TotalCards))]) % 4
	return cards, trump
}

// Send cards and information about round to each player.
// Att stands for the player, who bets and moves on the first board of round before others.
func sendCards(shuffled []datatypes.Card, N, tr, att int32, tunnels [3]datatypes.Tunnel, df int32) {
	// Fill cards for each player
	cards := make([][]datatypes.Card, 3)
	for i := 0; i < 3; i++ {
		cards[i] = make([]datatypes.Card, N)
	}
	for i := int32(0); i < N; i++ {
		cards[0][i] = shuffled[3*i]
		cards[1][i] = shuffled[3*i+1]
		cards[2][i] = shuffled[3*i+2]
	}

	// Send info
	r := datatypes.RoundInfo{NumOfCards: N, Cards: nil, Attacker: datatypes.Player(att), Trump: tr, DarkFlag: df}
	for i := 0; i < 3; i++ {
		r.Cards = cards[i]
		tunnels[i].NewRound <- r
	}
}

// Get bets from each player and inform others about it.
// Att stands for the first player to bet
func getBets(att int32, tunnels [3]datatypes.Tunnel) (bets [3]int32) {
	for i := int32(0); i < 3; i++ {
		bet := <-tunnels[(att+i)%3].BetToProc
		bets[bet.Player] = bet.Bet
		// Send bet to other players
		for j := 0; j < 3; j++ {
			tunnels[j].BetFromProc <- bet
		}
	}
	return
}

// holdBoard holds moves on the board and returns the winner of the board.
// Att stands for the player, who moves first on this board
func holdBoard(att, trump int32, tunnels [3]datatypes.Tunnel) (winner int32) {
	// Create a board to store cards
	board := [3]datatypes.Card{0, 0, 0}

	// Fill board
	for i := int32(0); i < 3; i++ {
		move := <-tunnels[(att+i)%3].MoveToProc
		board[i] = move.Card
		// Send move to all players
		for j := 0; j < 3; j++ {
			tunnels[j].MoveFromProc <- move
		}
	}

	// Calculate who wins the board
	winner = calculateWinner(board, trump, att)

	// Tell who wins
	for i := 0; i < 3; i++ {
		tunnels[i].WinnerName <- datatypes.WinInfo{Player: datatypes.Player(winner)}
	}

	return
}

// Calculate results of the round, considering bets and actual wins.
func calculateResults(bets, wins [3]int32) (res [3]int32) {
	for i := 0; i < 3; i++ {
		if wins[i] > bets[i] {
			res[i] = wins[i]
		} else if wins[i] == bets[i] {
			res[i] = bets[i] * 10
		} else {
			res[i] = -10 * bets[i]
		}
	}
	return
}

// Update the table with the rusults of the last round.
func updateTable(table, res [3]int32) [3]int32 {
	for i := 0; i < 3; i++ {
		table[i] += res[i]
	}
	return table
}

// sendRes sends points of each player to handlers.
func sendRes(res [3]int32, tunnels [3]datatypes.Tunnel) {
	for i := 0; i < 3; i++ {
		tunnels[i].RoundResults <- datatypes.ResultsInfo{Res: res}
	}
}

// calculateWinner returns the index of winner of the board.
func calculateWinner(b [3]datatypes.Card, tr, att int32) (winner int32) {
	if !beats(b[0], b[1], tr) && !beats(b[0], b[2], tr) {
		return att
	}
	if !beats(b[1], b[2], tr) {
		return (att + 1) % 3
	}
	return (att + 2) % 3
}

// beats returns true, if second cards beats first, considering trump.
func beats(a, b datatypes.Card, tr int32) bool {
	trA := int32(a) % 4
	trB := int32(b) % 4
	valA := int32(a) / 4
	valB := int32(b) / 4
	if (trA == trB && valB > valA) || (trB == tr && trA != tr) {
		return true
	}
	return false
}
