// Package handler is designed to connect players with "processor", and transort data between them.
// The callable functions are:
// func HandleFirstConn(net.Conn, datatypes.tunnel, chan struct{}, chan<- int32)
// func HandleConn     (net.Conn, datatypes.Player, datatypes.Tunnel, chan struct{}, int32)
package handler

import (
	"bufio"
	"net"

	"github.com/TblKVANYA/hussarGame/server/datatypes"
	"github.com/TblKVANYA/hussarGame/server/handler/clientconn"
	"github.com/TblKVANYA/hussarGame/server/handler/procconn"
)

// handles the connection with the first who joined until it becomes the same with others
func HandleFirstConn(conn net.Conn, tunnel datatypes.Tunnel, done chan struct{}, numberChan chan<- int32) {
	defer func() { done <- struct{}{} }()
	defer conn.Close()

	clientRW := initReadWriter(conn)

	clientconn.SendIndex(clientRW, 0)
	clientconn.GetHello(clientRW)
	N := clientconn.GetNumber(clientRW)

	numberChan <- N
	close(numberChan)
	handleConn(clientRW, 0, tunnel, N)
}

// handles the beginning of connection with everyone except the first who joined
func HandleConn(conn net.Conn, index datatypes.Player, tunnel datatypes.Tunnel, done chan struct{}, N int32) {
	defer func() { done <- struct{}{} }()

	defer conn.Close()

	clientRW := initReadWriter(conn)

	clientconn.SendIndex(clientRW, int32(index))
	clientconn.GetHello(clientRW)
	clientconn.SendNumberOfPlayers(clientRW, N)

	handleConn(clientRW, index, tunnel, N)
}

// inits a ReadWriter
func initReadWriter(conn net.Conn) *bufio.ReadWriter {
	clientReader := bufio.NewReader(conn)
	clientWriter := bufio.NewWriter(conn)
	clientRW := bufio.NewReadWriter(clientReader, clientWriter)
	return clientRW
}

// handles connection between player and server.
func handleConn(clientRW *bufio.ReadWriter, index datatypes.Player, tunnel datatypes.Tunnel, playersN int32) {

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
			provideBets(clientRW, tunnel, round, index, playersN)
		} else {
			provideBets(clientRW, tunnel, round, index, playersN)
			clientconn.SendDeck(clientRW, round)
		}

		// Cur is a player, who moves first on the next board
		cur := datatypes.WinInfo{Player: round.Attacker}
		N := round.NumOfCards
		// Card moves
		for i := int32(0); i < N; i++ {
			provideBoard(clientRW, tunnel, cur, index, playersN)
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
func provideBets(clientRW *bufio.ReadWriter, tunnel datatypes.Tunnel, round datatypes.RoundInfo, index datatypes.Player, playersN int32) {
	for i := int32(0); i < playersN; i++ {
		if (round.Attacker+datatypes.Player(i))%datatypes.Player(playersN) == index {
			b := clientconn.GetBet(clientRW)
			procconn.SendBet(tunnel, index, b)
		}
		bet := procconn.GetBet(tunnel)
		clientconn.SendBet(clientRW, bet)
	}
}

// provideBoard holds client-server data exchange during a board
func provideBoard(clientRW *bufio.ReadWriter, tunnel datatypes.Tunnel, att datatypes.WinInfo, index datatypes.Player, playersN int32) {
	for i := int32(0); i < playersN; i++ {
		if (att.Player+datatypes.Player(i))%datatypes.Player(playersN) == index {
			card := clientconn.GetCard(clientRW)
			procconn.SendMove(tunnel, index, card)
		}
		move := procconn.GetMove(tunnel)
		clientconn.SendCard(clientRW, move)
	}
}
