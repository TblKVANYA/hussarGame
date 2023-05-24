// Package processor implement the processor of the game, which encapsulates much of the logic of the program.
// The only callable function is Processor([]datatypes.Tunnel, int32)
package processor

import (
	"math/rand"

	"github.com/TblKVANYA/hussarGame/server/datatypes"
)

// holdRound holds the whole round with cardsN cards and playesN players and returns updated table of players ans scores.
// Attacker stands for the player, who bets and moves on the first board of round before others.
// Fl can take values of 0 (usual round), 1 (dark round), 2 (untrumped round) and 3(golden round)
func holdRound(attacker, cardsN int32, table []int32, random *rand.Rand, tunnels []datatypes.Tunnel, fl int, playersN int32) []int32 {
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
	sendCards(shuffled, cardsN, trump, attacker, tunnels, df, playersN)

	// Get bets
	bets := getBets(attacker, tunnels, playersN)

	// Hold boards
	wins := make([]int32, playersN)
	for n := int32(0); n < cardsN; n++ {
		attacker = holdBoard(attacker, trump, tunnels, playersN)
		wins[attacker]++
	}

	// Calculate results of round
	results := calculateResults(bets, wins, playersN)
	// Is it golden round?
	if fl == 3 {
		for i := int32(0); i < playersN; i++ {
			results[i] *= 3
		}
	}
	// End round
	sendRes(results, tunnels, playersN)
	table = updateTable(table, results, playersN)

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
func sendCards(shuffled []datatypes.Card, N, tr, att int32, tunnels []datatypes.Tunnel, df int32, playersN int32) {
	// Fill cards for each player
	cards := make([][]datatypes.Card, playersN)
	for i := int32(0); i < playersN; i++ {
		cards[i] = make([]datatypes.Card, N)
	}
	for i := int32(0); i < N; i++ {
		for j := int32(0); j < playersN; j++ {
			cards[j][i] = shuffled[playersN*i+j]
		}
	}

	// Send info
	r := datatypes.RoundInfo{NumOfCards: N, Cards: nil, Attacker: datatypes.Player(att), Trump: tr, DarkFlag: df}
	for i := int32(0); i < playersN; i++ {
		r.Cards = cards[i]
		tunnels[i].NewRound <- r
	}
}

// Get bets from each player and inform others about it.
// Att stands for the first player to bet
func getBets(att int32, tunnels []datatypes.Tunnel, playersN int32) (bets []int32) {
	bets = make([]int32, playersN)
	for i := int32(0); i < playersN; i++ {
		bet := <-tunnels[(att+i)%playersN].BetToProc
		bets[bet.Player] = bet.Bet
		// Send bet to other players
		for j := int32(0); j < playersN; j++ {
			tunnels[j].BetFromProc <- bet
		}
	}
	return
}

// holdBoard holds moves on the board and returns the winner of the board.
// Att stands for the player, who moves first on this board
func holdBoard(att, trump int32, tunnels []datatypes.Tunnel, playersN int32) (winner int32) {
	// Create a board to store cards
	board := make([]datatypes.Card, playersN)

	// Fill board
	for i := int32(0); i < playersN; i++ {
		move := <-tunnels[(att+i)%playersN].MoveToProc
		board[i] = move.Card
		// Send move to all players
		for j := int32(0); j < playersN; j++ {
			tunnels[j].MoveFromProc <- move
		}
	}

	// Calculate who wins the board
	winner = calculateWinner(board, trump, att, playersN)

	// Tell who wins
	for i := int32(0); i < playersN; i++ {
		tunnels[i].WinnerName <- datatypes.WinInfo{Player: datatypes.Player(winner)}
	}

	return
}

// Calculate results of the round, considering bets and actual wins.
func calculateResults(bets, wins []int32, playersN int32) (res []int32) {
	res = make([]int32, playersN)
	for i := int32(0); i < playersN; i++ {
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
func updateTable(table, res []int32, playersN int32) []int32 {
	for i := int32(0); i < playersN; i++ {
		table[i] += res[i]
	}
	return table
}

// sendRes sends points of each player to handlers.
func sendRes(res []int32, tunnels []datatypes.Tunnel, playersN int32) {
	for i := int32(0); i < playersN; i++ {
		tunnels[i].RoundResults <- datatypes.ResultsInfo{Res: res}
	}
}

// calculateWinner returns the index of winner of the board.
func calculateWinner(b []datatypes.Card, tr, att int32, playersN int32) (winner int32) {
	leader := datatypes.MoveInfo{Player: 0, Card: b[0]}

	for i := int32(1); i < playersN; i++ {
		if beats(leader.Card, b[i], tr) {
			leader = datatypes.MoveInfo{Player: datatypes.Player(i), Card: b[i]}
		}
	}

	return (att + int32(leader.Player)) % playersN
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
