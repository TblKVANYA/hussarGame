// package procconn holds the connection between "processor" and "handler".
// It is only used by the "handler" to avoid confusion between "handler"-"processor" and "handler"-client data exchange.
// Functions are:
// SendHello(datatypes.Tunnel, datatypes.Player)
// SendBet  (datatypes.Tunnel, datatypes.Player, int32)
// GetInfo  (datatypes.Tunnel) (datatypes.RoundInfo)
// GetBet   (datatypes.Tunnel) (datatypes.BetInfo)
// GetMove  (datatypes.Tunnel) (datatypes.MoveInfo)
// GetWinner(datatypes.Tunnel) (datatypes.WinInfo)
// GetRes   (datatypes.Tunnel) (datatypes.ResultsInfo)
package procconn

import "github.com/TblKVANYA/hussarGame/server/datatypes"

// SendHello tells the "processor" that player has joined the game.
func SendHello(tunnel datatypes.Tunnel, index datatypes.Player) {
	tunnel.MoveToProc <- datatypes.MoveInfo{Player: index, Card: datatypes.Card(0)}
}

// SendBet sends player's bet to the "processor".
func SendBet(tunnel datatypes.Tunnel, index datatypes.Player, b int32) {
	tunnel.BetToProc <- datatypes.BetInfo{Player: index, Bet: b}
}

// SendMove sends player's move to the "processor".
func SendMove(tunnel datatypes.Tunnel, index datatypes.Player, c datatypes.Card) {
	tunnel.MoveToProc <- datatypes.MoveInfo{Player: index, Card: c}
}

// GetInfo gets information about upcoming round.
func GetInfo(tunnel datatypes.Tunnel) (info datatypes.RoundInfo) {
	info = <-tunnel.NewRound
	return
}

// GetBet gets information about one's bet.
func GetBet(tunnel datatypes.Tunnel) (bet datatypes.BetInfo) {
	bet = <-tunnel.BetFromProc
	return
}

// GetMove gets information about one's move.
func GetMove(tunnel datatypes.Tunnel) (move datatypes.MoveInfo) {
	move = <-tunnel.MoveFromProc
	return
}

// GetWinner gets information who wins passed board
func GetWinner(tunnel datatypes.Tunnel) (winner datatypes.WinInfo) {
	winner = <-tunnel.WinnerName
	return
}

// GetRes gets results of passed round for each player
func GetRes(tunnel datatypes.Tunnel) (res datatypes.ResultsInfo) {
	res = <-tunnel.RoundResults
	return
}
