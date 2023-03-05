// package procconn holds the connection between client and "handler".
// It is designed to avoid confusion between "handler"-"processor" and "handler"-client data exchange.
// Functions are:
// func SendIndex          (net.Conn, int32)
// func SendNumberOfCards  (net.Conn, int32)
// func SendAttackerAndFlag(net.Conn, datatypes.Player, int32)
// func SendDeck           (net.Conn, datatypes.RoundInfo)
// func SendBet            (net.Conn, datatypes.BetInfo)
// func SendCard           (net.Conn, datatypes.MoveInfo)
// func SendWinner         (net.Conn, datatypes.WinInfo)
// func SendRes            (net.Conn, datatypes.ResultsInfo)
// func SendGB             (net.Conn)
// func GetHello           (net.Conn)
// func GetBet             (net.Conn) (int32)
// func GetCard            (net.Conn) (datatypes.Card)
package clientconn

import (
	"fmt"
	"log"
	"net"

	"github.com/TblKVANYA/hussarGame/server/datatypes"
)

// Sends player his index.
func SendIndex(conn net.Conn, index int32) {
	fmt.Fprintln(conn, index)
}

// Sends player number of cards in upcoming round.
func SendNumberOfCards(conn net.Conn, N int32) {
	fmt.Fprintln(conn, N)
}

// Sends major information about upcoming round.
func SendAttackerAndFlag(conn net.Conn, i datatypes.Player, fl int32) {
	fmt.Fprintln(conn, int32(i))
	fmt.Fprintln(conn, fl)
}

// Sends player's deck to him.
func SendDeck(conn net.Conn, info datatypes.RoundInfo) {
	for i := int32(0); i < info.NumOfCards; i++ {
		fmt.Fprintln(conn, int32(info.Cards[i]))
	}
	fmt.Fprintln(conn, int32(info.Trump))
}

// Sends one's bet to the player.
func SendBet(conn net.Conn, b datatypes.BetInfo) {
	fmt.Fprintln(conn, b.Bet)
}

// Sends one's move to the player.
func SendCard(conn net.Conn, c datatypes.MoveInfo) {
	fmt.Fprintln(conn, int32(c.Card))
}

// Sends who has won the board to the player.
func SendWinner(conn net.Conn, i datatypes.WinInfo) {
	fmt.Fprintln(conn, int32(i.Player))
}

// Sends results of the round to the player.
func SendRes(conn net.Conn, t datatypes.ResultsInfo) {
	fmt.Fprintf(conn, "%d %d %d\n", t.Res[0], t.Res[1], t.Res[2])
}

// Sends flag of game ending to the player.
func SendGB(conn net.Conn) {
	fmt.Fprintln(conn, -1)
}

// Gets smth from to player to ensure he is ready.
func GetHello(conn net.Conn) {
	s := ""
	_, err := fmt.Fscan(conn, &s)
	if err != nil {
		log.Fatal(err)
	}
}

// Gets a bet from the player.
func GetBet(conn net.Conn) (bet int32) {
	_, err := fmt.Fscan(conn, &bet)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// Gets a card from the player.
func GetCard(conn net.Conn) (c datatypes.Card) {
	var val int32
	_, err := fmt.Fscan(conn, &val)
	if err != nil {
		log.Fatal(err)
	}
	return datatypes.Card(val)
}
