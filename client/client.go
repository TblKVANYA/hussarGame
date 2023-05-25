// package client contains a client for Hussar card game, which can be started with the function
// func Client(string)
package client

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/TblKVANYA/hussarGame/client/drawer"
)

// func Client starts a client for Hussar
func Client(ip string) {
	if net.ParseIP(ip) == nil {
		log.Fatal("wrong ip")
	}

	// establish connection
	addr := ip + ":8088"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	serverReader := bufio.NewReader(conn)
	serverWriter := bufio.NewWriter(conn)
	serverRW := bufio.NewReadWriter(serverReader, serverWriter)

	myIndex := getIndex(serverRW)

	sendHello(serverRW)

	var playersN int32
	if myIndex == 0 {
		playersN = sendNumberOfPlayers(serverRW)
	} else {
		playersN = getNumberOfPlayers(serverRW)
	}
	drawer.Init(playersN, myIndex)

	// Not to let player bet "0" playersN times in a row
	myBets := make([]int32, 0)

	for {
		cardsN := getNumber(serverRW)

		// End flag
		if cardsN == -1 {
			break
		}
		// Scan primary info
		att, darkFlag := getInfo(serverRW)

		drawer.UpdateIndex(cardsN)
		var myCards []int32

		// Dark round check
		if darkFlag == 0 {
			myCards = getCards(serverRW, cardsN)
			myBets = takeBets(serverRW, myBets, att, myIndex, playersN)
		} else {
			myBets = takeBets(serverRW, myBets, att, myIndex, playersN)
			myCards = getCards(serverRW, cardsN)
		}

		// Holding a board and its results
		for n := cardsN; n > 0; n-- {
			takeBoard(serverRW, myCards, n, att, myIndex, playersN)
			drawer.UpdateIndex(n - 1)

			att = getWinner(serverRW)
			drawer.AddWin(att)
			drawer.FlushBoard()
		}

		// Get results of the round
		tmp := int32(0)
		for i := int32(0); i < playersN; i++ {
			tmp = getRes(serverRW)
			drawer.AddRes(i, tmp)
		}
		drawer.FlushRound()
	}
}

// takeBets leads player through bets
func takeBets(rw *bufio.ReadWriter, myBets []int32, att, myIndex, playersN int32) []int32 {
	var arrState int32
	for i := int32(0); i < playersN; i++ {
		arrState = (playersN + att + i - myIndex) % playersN
		drawer.SetArrow(arrState)
		drawer.Draw()
		var b int32
		// time for own bet
		if arrState == 0 {
			b, myBets = scanUserBet(myBets, playersN)
			sendBet(rw, b)
		}

		b = getBet(rw)
		drawer.SetBet((att+i)%playersN, b)
	}
	return myBets
}

// takeBoard leads client through board
func takeBoard(rw *bufio.ReadWriter, myCards []int32, totalCards, att, myIndex, playersN int32) {
	var arrState int32
	for i := int32(0); i < playersN; i++ {
		arrState = (playersN + att + i - myIndex) % playersN
		drawer.SetArrow(arrState)
		drawer.Draw()
		var c int32

		// Time for own move
		if arrState == 0 {
			c = scanCardIndex(totalCards)
			drawer.ReplaceCard(c-1, totalCards-1)
			myCards[c-1], myCards[totalCards-1] = myCards[totalCards-1], myCards[c-1]
			if myCards[totalCards-1] != drawer.GetMagic() {
				sendCard(rw, myCards[totalCards-1])
			} else {
				sendCard(rw, parseMagic())
			}
		}

		// store visible changes of the board
		drawer.FlushCard((att+i)%playersN, totalCards-1)
		c = getCard(rw)
		drawer.SetCard((att+i)%playersN, c)
	}
	drawer.Draw()

	fmt.Println("Press Enter to continue")
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
}

// sendHello makes sure player is ready and send hello to server
func sendHello(rw *bufio.ReadWriter) {
	// Some rules
	fmt.Println("That's a game with hard rules. I'll explain them.")
	fmt.Println("Remember, that J♧ is the magic card. It could be interpretated as any of 36 usual cards. Moreover, it could be interpretated as a 5 of any suit")
	fmt.Println("What's more, it could be a S(uperior) of any suit. S is a rank superior to any other")
	fmt.Println("To choose your card print its index(written under the card)")
	fmt.Println("Do not send 0 as a bet for many times in a row.")
	fmt.Println("There is a dark round with 12 cards per player, where you bet before watching cards.")
	fmt.Println("Also there is an untrumped round.")
	fmt.Println("Points in the last round will be multiplyed by 3.")

	// Some greeting
	fmt.Println("Print smth (preferably \"go\") to start")
	s := ""
	for s == "" {
		fmt.Scanf("%s", &s)
	}
	err := writeInt32(rw.Writer, 0)
	if err != nil {
		log.Fatal(err.Error() + " in sendHello")
	}
	err = rw.Writer.Flush()
	if err != nil {
		log.Fatal(err.Error() + " in sendHello")
	}
	fmt.Println("Waiting for others")
}

// sendNumberOfPlayers gives a choice how many players will be playing for a user and sends his option to server.
func sendNumberOfPlayers(rw *bufio.ReadWriter) int32 {
	n := scanNumberOfPlayers()
	err := writeInt32(rw.Writer, n)
	if err != nil {
		log.Fatal(err.Error() + " in sendNumberOfPlayers")
	}
	err = rw.Writer.Flush()
	if err != nil {
		log.Fatal(err.Error() + " in sendNumberOfPlayers")
	}
	return n
}

// sendBet sends bet b to rw
func sendBet(rw *bufio.ReadWriter, b int32) {
	err := writeInt32(rw.Writer, b)
	if err != nil {
		log.Fatal(err.Error() + " in sendBet")
	}
	err = rw.Writer.Flush()
	if err != nil {
		log.Fatal(err.Error() + " in sendBet")
	}
}

// sendCard sends card c to rw
func sendCard(rw *bufio.ReadWriter, c int32) {
	err := writeInt32(rw.Writer, c)
	if err != nil {
		log.Fatal(err.Error() + " in sendCard")
	}
	err = rw.Writer.Flush()
	if err != nil {
		log.Fatal(err.Error() + " in sendCard")
	}
}

// Writes int32 to w
func writeInt32(w *bufio.Writer, val int32) error {
	return binary.Write(w, binary.BigEndian, val)
}

// getIndex gets client's index from rw.
func getIndex(rw *bufio.ReadWriter) (i int32) {
	i, err := readInt32(rw.Reader)
	if err != nil {
		log.Fatal(err.Error() + " in getIndex")
	}
	return
}

// getNumberOfPlayers gets number of participants from rw.
func getNumberOfPlayers(rw *bufio.ReadWriter) (n int32) {
	n, err := readInt32(rw.Reader)
	if err != nil {
		log.Fatal(err.Error() + " in getPN")
	}
	return n
}

// getInfo gets majot info about upcoming round
func getInfo(rw *bufio.ReadWriter) (att, df int32) {
	vals, err := readFewInt32(rw.Reader, 2)
	if err != nil {
		log.Fatal(err.Error() + " in getInfo")
	}
	_ = vals[1]
	return vals[0], vals[1]
}

// getNumber gets a number of cards in upcoming round
func getNumber(rw *bufio.ReadWriter) (n int32) {
	n, err := readInt32(rw.Reader)
	if err != nil {
		log.Fatal(err.Error() + " in getNumber")
	}
	return
}

// getCards gets N cards from rw and returns slice of these cards
func getCards(rw *bufio.ReadWriter, N int32) []int32 {
	cards := make([]int32, N)
	for i := int32(0); i < N; i++ {
		cards[i] = getCard(rw)
	}

	// Find out  trump
	tr := getTrump(rw)

	drawer.SetDecks(cards)
	drawer.SetTrump(tr)
	return cards
}

// getTrump gets a trump from rw.
func getTrump(rw *bufio.ReadWriter) (tr int32) {
	tr, err := readInt32(rw.Reader)
	if err != nil {
		log.Fatal(err.Error() + " in getTrump")
	}
	return
}

// getBet gets one player bet from rw.
func getBet(rw *bufio.ReadWriter) (b int32) {
	b, err := readInt32(rw.Reader)
	if err != nil {
		log.Fatal(err.Error() + " in getBet")
	}
	return
}

// getCard gets a card from rw.
func getCard(rw *bufio.ReadWriter) (c int32) {
	c, err := readInt32(rw.Reader)
	if err != nil {
		log.Fatal(err.Error() + " in getCard")
	}
	return
}

// getWinner gets winner index from rw.
func getWinner(rw *bufio.ReadWriter) (w int32) {
	w, err := readInt32(rw.Reader)
	if err != nil {
		log.Fatal(err.Error() + " in getWinner")
	}
	return
}

// getRes gets one player result score for a round from rw.
func getRes(rw *bufio.ReadWriter) (r int32) {
	r, err := readInt32(rw.Reader)
	if err != nil {
		log.Fatal(err.Error() + " in getRes")
	}
	return
}

// Reads n int32 values from r, not divided by anything.
func readFewInt32(r *bufio.Reader, n int) ([]int32, error) {
	p := make([]byte, 4)
	res := make([]int32, n)
	for i := 0; i < n; i++ {
		b, err := r.Read(p)
		if err != nil {
			return nil, err
		}
		if b != 4 {
			return nil, errors.New("read: number of read bytes isn't 4")
		}
		res[i] = int32(binary.BigEndian.Uint32(p))
	}
	return res, nil
}

// Reads int32 from r.
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

// scanNumberOfPlayers scans first player's choice of how many players will be playing.
func scanNumberOfPlayers() int32 {
	var n int32
	for n < 2 || n > 4 {
		fmt.Println("You have started a game. Choose how many players will play. Current options are: 2-4")
		fmt.Scanf("%d", &n)
	}
	return n
}

// scanUserBet scans player's bet from keyboard, checks it's correctness and sends it to server.
func scanUserBet(prevBets []int32, playersN int32) (int32, []int32) {
	var b int32 = -1
	for {
		fmt.Println("Send bet")
		fmt.Scanf("%d", &b)
		if b >= 0 && b <= 20 {
			if len(prevBets) < int(playersN)-1 {
				prevBets = append(prevBets, b)
				break
			}
			sum := b
			for i := len(prevBets) - int(playersN) + 1; i < len(prevBets); i++ {
				sum += prevBets[i]
			}
			if sum != 0 {
				prevBets = append(prevBets, b)
				break
			}
		}
	}
	return b, prevBets
}

// scanCardIndex is a function, which returns number of the chosen card from the slice of cards.
// n stands for number of cards in slice.
func scanCardIndex(n int32) (c int32) {
	for {
		fmt.Println("Choose card")
		fmt.Scanf("%d", &c)
		if c > 0 && c <= n {
			break
		}
	}
	return
}

// parseMagic is a function, which defines how player wants to use his J♧ card.
func parseMagic() int32 {
	fmt.Println("Please, choose how you want to use this card.")
	c, s := int32(0), ""

	// Parse rank
	for {
		fmt.Println("Write 5-10 or J, Q, K, A, S to choose rank.")
		fmt.Scanf("%s", &s)
		num, err := strconv.Atoi(s)
		if err == nil {
			if num >= 5 && num <= 10 {
				c += (int32(num) - 5) * 4
				break
			}
		} else {
			if s == "J" || s == "j" {
				c += 24
				break
			} else if s == "Q" || s == "q" {
				c += 28
				break
			} else if s == "K" || s == "k" {
				c += 32
				break
			} else if s == "A" || s == "a" {
				c += 36
				break
			} else if s == "S" || s == "s" {
				c += 40
				break
			}
		}
	}

	// Parse suit
	for {
		fmt.Println("To choose suit write H(for ♡), S(for ♤), D(for ♢) or C(for ♧)")
		fmt.Scanf("%s", &s)
		if s == "H" || s == "h" {
			break
		} else if s == "D" || s == "d" {
			c += 1
			break
		} else if s == "S" || s == "s" {
			c += 2
			break
		} else if s == "C" || s == "c" {
			c += 3
			break
		}

	}

	return c
}
