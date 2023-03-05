// Package handler is designed to connect players with "processor", and transort data between them.
// The only callable function is HandleConn(net.Conn, datatypes.Player, datatypes.Tunnel, chan struct{}).
package handler

import (
	"net"

	"github.com/TblKVANYA/hussarGame/server/datatypes"
	"github.com/TblKVANYA/hussarGame/server/handler/clientconn"
	"github.com/TblKVANYA/hussarGame/server/handler/procconn"
)

// HandleConn handles connection between player and server.
func HandleConn(conn net.Conn, index datatypes.Player, tunnel datatypes.Tunnel, done chan struct{}) {
	defer func() { done <- struct{}{} }()

	// Sensual foreplay
	clientconn.SendIndex(conn, int32(index))
	clientconn.GetHello(conn)
	procconn.SendHello(tunnel, index)

	for {
		round := procconn.GetInfo(tunnel)
		// End signal
		if round.NumOfCards == -1 {
			break
		}
		// Send must-have info to client
		clientconn.SendNumberOfCards(conn, round.NumOfCards)
		clientconn.SendAttackerAndFlag(conn, round.Attacker, round.DarkFlag)

		// Dark round swaps order of client sending bets and watching his cards
		if round.DarkFlag == 0 {
			clientconn.SendDeck(conn, round)
			provideBets(conn, tunnel, round, index)
		} else {
			provideBets(conn, tunnel, round, index)
			clientconn.SendDeck(conn, round)
		}

		// Cur is a player, who moves first on the next board
		cur := datatypes.WinInfo{Player: round.Attacker}
		N := round.NumOfCards
		// Card moves
		for i := int32(0); i < N; i++ {
			provideBoard(conn, tunnel, cur, index)
			cur = procconn.GetWinner(tunnel)
			clientconn.SendWinner(conn, cur)
		}
		res := procconn.GetRes(tunnel)
		clientconn.SendRes(conn, res)
	}
	clientconn.SendGB(conn)

	conn.Close()
	tunnel.CloseFromHandler()
}

// provideBets holds client-server data exchange during the time of bets.
func provideBets(conn net.Conn, tunnel datatypes.Tunnel, round datatypes.RoundInfo, index datatypes.Player) {
	for i := 0; i < 3; i++ {
		if (round.Attacker+datatypes.Player(i))%3 == index {
			b := clientconn.GetBet(conn)
			procconn.SendBet(tunnel, index, b)
		}
		bet := procconn.GetBet(tunnel)
		clientconn.SendBet(conn, bet)
	}
}

// provideBoard holds client-server data exchange during a board
func provideBoard(conn net.Conn, tunnel datatypes.Tunnel, att datatypes.WinInfo, index datatypes.Player) {
	for i := 0; i < 3; i++ {
		if (att.Player+datatypes.Player(i))%3 == index {
			card := clientconn.GetCard(conn)
			procconn.SendMove(tunnel, index, card)
		}
		move := procconn.GetMove(tunnel)
		clientconn.SendCard(conn, move)
	}
}
