// package client contains a client for Hussar card game, which can be started with the function
// func Client(string)
package client

import (
	"bufio"
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

	addr := ip + ":8088"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	var myIndex int32
	fmt.Fscanf(conn, "%d", &myIndex)
	drawer.SetSides(myIndex)

	// Not to let player bet "0" 3 times in a row
	myBets := make([]int32, 0)

	sendHello(conn)

	for {
		var N int32
		fmt.Fscanf(conn, "%d", &N)

		// End flag
		if N == -1 {
			break
		}
		// Scan primary info
		att, darkFlag := int32(0), 0
		fmt.Fscanf(conn, "%d", &att)
		fmt.Fscanf(conn, "%d", &darkFlag)

		drawer.UpdateIndex(N)
		var myCards []int32

		// Dark round check
		if darkFlag == 0 {
			myCards = gatherCards(conn, N)
			myBets = takeBets(conn, myBets, att, myIndex)
		} else {
			myBets = takeBets(conn, myBets, att, myIndex)
			myCards = gatherCards(conn, N)
		}

		// Holding a board and its results
		for n := N; n > 0; n-- {
			takeBoard(conn, myCards, n, att, myIndex)
			drawer.UpdateIndex(n - 1)

			fmt.Fscanf(conn, "%d", &att)
			drawer.AddWin(att)
			drawer.FlushBoard()
		}

		// Get results of the round
		tmp := int32(0)
		for i := int32(0); i < 3; i++ {
			fmt.Fscanf(conn, "%d", &tmp)
			drawer.AddRes(i, tmp)
		}
		drawer.FlushRound()
	}
}

// sendHello makes sure player is ready and send hello to server
func sendHello(conn net.Conn) {
	// Some rules
	fmt.Println("That's a game with hard rules. I'll explain them.")
	fmt.Println("Remember, that J♧ is the magic card. It could be interpretated as any of 36 usual cards. Moreover, it could be interpretated as a 5 of any suit")
	fmt.Println("What's more, it could be a S(uperior) of any suit. S is a rank superior to any other")
	fmt.Println("To choose your card print its index(written under the card)")
	fmt.Println("Do not send bet==0 three times in a row.")
	fmt.Println("There is a dark round with 12 cards per player, where you bet before watching cards")
	fmt.Println("Also there is an untrumped round. Points in the last round will be multiplyed by 3")

	// Some greeting
	fmt.Println("Print smth (preferably \"go\") to start")
	s := ""
	for s == "" {
		fmt.Scanf("%s", &s)
	}
	fmt.Fprintln(conn, s)
	fmt.Println("Waiting for others")
}

// gatherCards scans N cards from connection conn and returns slice of these cards
func gatherCards(conn net.Conn, N int32) []int32 {
	cards := make([]int32, N)
	for i := int32(0); i < N; i++ {
		fmt.Fscanf(conn, "%d", &(cards[i]))
	}
	// Find out  trump
	var tr int32
	fmt.Fscanf(conn, "%d", &tr)

	drawer.SetDecks(cards)
	drawer.SetTrump(tr)
	return cards
}

// takeBets leads player through bets
func takeBets(conn net.Conn, myBets []int32, att, myIndex int32) []int32 {
	var arrState int32
	for i := int32(0); i < 3; i++ {
		arrState = (3 + att + i - myIndex) % 3
		drawer.SetArrow(arrState)
		drawer.Draw()
		var b int32
		// time for own bet
		if arrState == 0 {
			b, myBets = scanBet(myBets)
			fmt.Fprintln(conn, b)
		}

		fmt.Fscanf(conn, "%d", &b)
		drawer.SetBet((att+i)%3, b)
	}
	return myBets
}

// scanBet scans player's bet, assume it's correct and sends it to server
func scanBet(prevBets []int32) (int32, []int32) {
	var b int32 = -1
	for {
		fmt.Println("Send bet")
		fmt.Scanf("%d", &b)
		if b >= 0 && b <= 12 {
			if len(prevBets) < 2 {
				prevBets = append(prevBets, b)
				break
			}
			if b != 0 || prevBets[len(prevBets)-1] != 0 || prevBets[len(prevBets)-2] != 0 {
				fmt.Println(b, prevBets[len(prevBets)-1], prevBets[len(prevBets)-2])
				prevBets = append(prevBets, b)
				break
			}
		}
	}
	return b, prevBets
}

// takeBoard leads client through board
func takeBoard(conn net.Conn, myCards []int32, totalCards, att, myIndex int32) {
	var arrState int32
	for i := int32(0); i < 3; i++ {
		arrState = (3 + att + i - myIndex) % 3
		drawer.SetArrow(arrState)
		drawer.Draw()
		var c int32

		// Time for own move
		if arrState == 0 {
			c = scanCardIndex(totalCards)
			drawer.ReplaceCard(c-1, totalCards-1)
			myCards[c-1], myCards[totalCards-1] = myCards[totalCards-1], myCards[c-1]
			if myCards[totalCards-1] != drawer.GetMagic() {
				fmt.Fprintln(conn, myCards[totalCards-1])
			} else {
				fmt.Fprintln(conn, parseMagic())
			}
		}

		// store visible changes of the board
		drawer.FlushCard((att+i)%3, totalCards-1)
		fmt.Fscanf(conn, "%d", &c)
		drawer.SetCard((att+i)%3, c)
	}
	drawer.Draw()

	fmt.Println("Press Enter to continue")
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
}

// scanCardIndex is a function, which returns number of the chosen card from the slice of cards.
// n stands for number of cards in slice
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
			c += (int32(num) - 5) * 4
			break
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
