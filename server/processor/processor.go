// Package processor implement the processor of the game, which encapsulates much of the logic of the program.
// The only callable function is Processor([]datatypes.Tunnel, int32)
package processor

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/TblKVANYA/hussarGame/server/datatypes"
)

// Processor encapsulates logic of the programme.
func Processor(tunnels []datatypes.Tunnel, playersN int32) {
	// Wait for connections from all the players
	for i := int32(0); i < playersN; i++ {
		<-tunnels[i].MoveToProc
		fmt.Println("Player has joined!")
	}

	time.Sleep(500 * time.Millisecond)
	table := make([]int32, playersN)

	// Init random
	source := rand.NewSource(time.Now().Unix())
	random := rand.New(source)
	attacker := random.Int31n(int32(playersN))

	// Hold usual rounds
	rules, ok := datatypes.UsualRulesArray[playersN]
	if !ok {
		log.Fatal("unsupported number of players")
	}
	for _, cardsN := range rules {
		table = holdRound(attacker, cardsN, table, random, tunnels, 0, playersN)
		attacker = (attacker + 1) % playersN
	}

	// Hold some extra rounds with unusual rules
	for i := 1; i <= 3; i++ {
		table = holdRound(attacker, datatypes.TotalCards/int32(playersN), table, random, tunnels, i, playersN)
		attacker = (attacker + 1) % playersN
	}

	// Send end signal to handlers and close chans
	endInfo := datatypes.RoundInfo{NumOfCards: -1, Cards: nil, Attacker: datatypes.Player(0), Trump: 0, DarkFlag: 0}
	for i := int32(0); i < playersN; i++ {
		tunnels[i].NewRound <- endInfo
		tunnels[i].CloseFromProcessor()
	}
}
