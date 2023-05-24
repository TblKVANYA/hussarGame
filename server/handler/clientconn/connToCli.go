// package procconn holds the connection between client and "handler".
// It is designed to avoid confusion between "handler"-"processor" and "handler"-client data exchange.
// Functions are:
// func SendNumberOfPlayers(*bufio.ReadWriter, int32)
// func SendIndex          (*bufio.ReadWriter, int32)
// func SendNumberOfCards  (*bufio.ReadWriter, int32)
// func SendAttackerAndFlag(*bufio.ReadWriter, datatypes.Player, int32)
// func SendDeck           (*bufio.ReadWriter, datatypes.RoundInfo)
// func SendBet            (*bufio.ReadWriter, datatypes.BetInfo)
// func SendCard           (*bufio.ReadWriter, datatypes.MoveInfo)
// func SendWinner         (*bufio.ReadWriter, datatypes.WinInfo)
// func SendRes            (*bufio.ReadWriter, datatypes.ResultsInfo)
// func SendGB             (*bufio.ReadWriter)
// func GetHello           (*bufio.ReadWriter)
// func GetNumber          (*bufio.ReadWriter) (int32)
// func GetBet             (*bufio.ReadWriter) (int32)
// func GetCard            (*bufio.ReadWriter) (datatypes.Card)
package clientconn

import (
	"bufio"
	"encoding/binary"
	"errors"
	"log"

	"github.com/TblKVANYA/hussarGame/server/datatypes"
)

// Sends player number of players.
func SendNumberOfPlayers(rw *bufio.ReadWriter, N int32) {
	err := writeInt32(rw.Writer, N)
	if err != nil {
		log.Fatal(err.Error() + " in SendNumberOfPlayers")
	}
	err = rw.Writer.Flush()
	if err != nil {
		log.Fatal(err.Error() + " in SendNumberOfPlayers")
	}
}

// Sends player his index.
func SendIndex(rw *bufio.ReadWriter, index int32) {
	err := writeInt32(rw.Writer, index)
	if err != nil {
		log.Fatal(err.Error() + " in SendIndex")
	}
	err = rw.Writer.Flush()
	if err != nil {
		log.Fatal(err.Error() + " in SendIndex")
	}
}

// Sends player number of cards in upcoming round.
func SendNumberOfCards(rw *bufio.ReadWriter, N int32) {
	err := writeInt32(rw.Writer, N)
	if err != nil {
		log.Fatal(err.Error() + " in SendNumberOfCards")
	}
	err = rw.Writer.Flush()
	if err != nil {
		log.Fatal(err.Error() + " in SendNumberOfCards")
	}
}

// Sends major information about upcoming round.
func SendAttackerAndFlag(rw *bufio.ReadWriter, i datatypes.Player, fl int32) {
	err := writeInt32(rw.Writer, int32(i))
	if err != nil {
		log.Fatal(err.Error() + " in SendAttackerAndFlag")
	}
	err = writeInt32(rw.Writer, fl)
	if err != nil {
		log.Fatal(err.Error() + " in SendAttackerAndFlag")
	}
	err = rw.Writer.Flush()
	if err != nil {
		log.Fatal(err.Error() + " in SendAttackerAndFlag")
	}
}

// Sends player's deck and trump.
func SendDeck(rw *bufio.ReadWriter, info datatypes.RoundInfo) {
	var err error
	for i := int32(0); i < info.NumOfCards; i++ {
		err = writeInt32(rw.Writer, int32(info.Cards[i]))
		if err != nil {
			log.Fatal(err.Error() + " in SendDeck")
		}
	}

	err = writeInt32(rw.Writer, int32(info.Trump))
	if err != nil {
		log.Fatal(err.Error() + " in SendDeck")
	}
	err = rw.Writer.Flush()
	if err != nil {
		log.Fatal(err.Error() + " in SendDeck")
	}
}

// Sends one's bet to the player.
func SendBet(rw *bufio.ReadWriter, b datatypes.BetInfo) {
	err := writeInt32(rw.Writer, b.Bet)
	if err != nil {
		log.Fatal(err.Error() + " in SendBet")
	}
	err = rw.Writer.Flush()
	if err != nil {
		log.Fatal(err.Error() + " in SendBet")
	}
}

// Sends one's move to the player.
func SendCard(rw *bufio.ReadWriter, c datatypes.MoveInfo) {
	err := writeInt32(rw.Writer, int32(c.Card))
	if err != nil {
		log.Fatal(err.Error() + " in SendCard")
	}
	err = rw.Writer.Flush()
	if err != nil {
		log.Fatal(err.Error() + " in SendCard")
	}
}

// Sends who has won the board to the player.
func SendWinner(rw *bufio.ReadWriter, i datatypes.WinInfo) {
	err := writeInt32(rw.Writer, int32(i.Player))
	if err != nil {
		log.Fatal(err.Error() + " in SendWinner")
	}
	err = rw.Writer.Flush()
	if err != nil {
		log.Fatal(err.Error() + " in SendWinner")
	}
}

// Sends results of the round to the player.
func SendRes(rw *bufio.ReadWriter, t datatypes.ResultsInfo) {
	var err error
	for i := range t.Res {
		err = writeInt32(rw.Writer, t.Res[i])
		if err != nil {
			log.Fatal(err.Error() + "in SendRes")
		}
	}
	err = rw.Writer.Flush()
	if err != nil {
		log.Fatal(err.Error() + " in SendRes")
	}
}

// Sends flag of game ending to the player.
func SendGB(rw *bufio.ReadWriter) {
	err := writeInt32(rw.Writer, int32(-1))
	if err != nil {
		log.Fatal(err.Error() + " in SendGb")
	}
	err = rw.Writer.Flush()
	if err != nil {
		log.Fatal(err.Error() + " in SendGb")
	}
}

// Gets smth from to player to ensure he is ready.
func GetHello(rw *bufio.ReadWriter) {
	_, err := readInt32(rw.Reader)
	if err != nil {
		log.Fatal(err)
	}
}

// Gets number of players from the first one to join.
func GetNumber(rw *bufio.ReadWriter) int32 {
	n, err := readInt32(rw.Reader)
	if err != nil {
		log.Fatal(err.Error() + " in GetNumber")
	}
	return n
}

// Gets a bet from the player.
func GetBet(rw *bufio.ReadWriter) int32 {
	b, err := readInt32(rw.Reader)
	if err != nil {
		log.Fatal(err.Error() + " in GetBet")
	}
	return b
}

// Gets a card from the player.
func GetCard(rw *bufio.ReadWriter) datatypes.Card {
	card, err := readInt32(rw.Reader)
	if err != nil {
		log.Fatal(err.Error() + " in GetCard")
	}
	return datatypes.Card(card)
}

// Writes int32 value to w.
func writeInt32(w *bufio.Writer, val int32) error {
	return binary.Write(w, binary.BigEndian, val)
}

// Reads int32 value from r.
func readInt32(r *bufio.Reader) (int32, error) {
	p := make([]byte, 4)
	n, err := r.Read(p)
	if err != nil {
		return 0, err
	}
	if n != 4 {
		return 0, errors.New("read: number of read bytes isn't 4")
	}

	val := int32(binary.BigEndian.Uint32(p))
	return val, nil
}
