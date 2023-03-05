// Package processor implement the processor of the game, which encapsulates much of the logic of the program.
// The only callable function is Processor([3]datatypes.Tunnel)
package processor

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/TblKVANYA/hussarGame/server/datatypes"
)

// Processor encapsulates logic of the programme.
func Processor(tunnels [3]datatypes.Tunnel) {
	// Wait for connections from all the players
	for i := 0; i < 3; i++ {
		<-tunnels[i].MoveToProc
		fmt.Println("Player has joined!")
	}
	time.Sleep(500 * time.Millisecond)
	table := [3]int32{0, 0, 0}

	// Init random
	source := rand.NewSource(time.Now().Unix())
	random := rand.New(source)
	attacker := random.Int31n(3)

	// Hold usual rounds
	for _, cardsN := range datatypes.UsualRulesArray {
		table = holdRound(attacker, cardsN, table, random, tunnels, 0)
		attacker = (attacker + 1) % 3
	}

	// Hold some extra rounds with unusual rules
	for i := 1; i <= 3; i++ {
		table = holdRound(attacker, 12, table, random, tunnels, i)
		attacker = (attacker + 1) % 3
	}

	// Send end signal to handlers and close chans
	endInfo := datatypes.RoundInfo{NumOfCards: -1, Cards: nil, Attacker: datatypes.Player(0), Trump: 0, DarkFlag: 0}
	for i := 0; i < 3; i++ {
		tunnels[i].NewRound <- endInfo
		tunnels[i].CloseFromProcessor()
	}
}
