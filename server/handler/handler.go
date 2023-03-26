// Package handler is designed to connect players with "processor", and transort data between them.
// The only callable function is HandleConn(net.Conn, datatypes.Player, datatypes.Tunnel, chan struct{}).
package handler

import (
	"bufio"
	"net"

	"github.com/TblKVANYA/hussarGame/server/datatypes"
	"github.com/TblKVANYA/hussarGame/server/handler/clientconn"
	"github.com/TblKVANYA/hussarGame/server/handler/procconn"
)

// HandleConn handles connection between player and server.
func HandleConn(conn net.Conn, index datatypes.Player, tunnel datatypes.Tunnel, done chan struct{}) {
	defer func() { done <- struct{}{} }()

	defer conn.Close()

	clientReader := bufio.NewReader(conn)
	clientWriter := bufio.NewWriter(conn)
	clientRW := bufio.NewReadWriter(clientReader, clientWriter)

	// Sensual foreplay
	clientconn.SendIndex(clientRW, int32(index))
	clientconn.GetHello(clientRW)
	procconn.SendHello(tunnel, index)

	for {
		round := procconn.GetInfo(tunnel)
		// End signal
		if round.NumOfCards == -1 {
			break
		}
		// Send must-have info to client
		clientconn.SendNumberOfCards(clientRW, round.NumOfCards)
		clientconn.SendAttackerAndFlag(clientRW, round.Attacker, round.DarkFlag)

		// Dark round swaps order of client sending bets and watching his cards
		if round.DarkFlag == 0 {
			clientconn.SendDeck(clientRW, round)
			provideBets(clientRW, tunnel, round, index)
		} else {
			provideBets(clientRW, tunnel, round, index)
			clientconn.SendDeck(clientRW, round)
		}

		// Cur is a player, who moves first on the next board
		cur := datatypes.WinInfo{Player: round.Attacker}
		N := round.NumOfCards
		// Card moves
		for i := int32(0); i < N; i++ {
			provideBoard(clientRW, tunnel, cur, index)
			cur = procconn.GetWinner(tunnel)
			clientconn.SendWinner(clientRW, cur)
		}
		res := procconn.GetRes(tunnel)
		clientconn.SendRes(clientRW, res)
	}
	clientconn.SendGB(clientRW)

	tunnel.CloseFromHandler()
}

// provideBets holds client-server data exchange during the time of bets.
func provideBets(clientRW *bufio.ReadWriter, tunnel datatypes.Tunnel, round datatypes.RoundInfo, index datatypes.Player) {
	for i := 0; i < 3; i++ {
		if (round.Attacker+datatypes.Player(i))%3 == index {
			b := clientconn.GetBet(clientRW)
			procconn.SendBet(tunnel, index, b)
		}
		bet := procconn.GetBet(tunnel)
		clientconn.SendBet(clientRW, bet)
	}
}

// provideBoard holds client-server data exchange during a board
func provideBoard(clientRW *bufio.ReadWriter, tunnel datatypes.Tunnel, att datatypes.WinInfo, index datatypes.Player) {
	for i := 0; i < 3; i++ {
		if (att.Player+datatypes.Player(i))%3 == index {
			card := clientconn.GetCard(clientRW)
			procconn.SendMove(tunnel, index, card)
		}
		move := procconn.GetMove(tunnel)
		clientconn.SendCard(clientRW, move)
	}
}
