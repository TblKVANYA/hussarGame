// datatypes provides the server with data types and structures
package datatypes

const TotalCards int32 = 36

// Hepls to differ card from other int values to avoid some mistakes
type Card int32

// Helps to differ player index from other int values to aviod some mistakes
type Player int32

// Information about one's move: who made it and which cards he used
type MoveInfo struct {
	Player
	Card
}

// Information about one's bet: who made it and what exactly it is
type BetInfo struct {
	Player
	Bet int32
}

// Information about upcoming round
type RoundInfo struct {
	NumOfCards int32
	Cards      []Card
	Attacker   Player
	Trump      int32
	DarkFlag   int32
}

// Information who has won the past board
type WinInfo struct {
	Player
}

// Information about results of the past round
type ResultsInfo struct {
	Res [3]int32
}

// A union of chans, connecting "processor" with "handler"
type Tunnel struct {
	MoveToProc   chan MoveInfo
	MoveFromProc chan MoveInfo
	BetToProc    chan BetInfo
	BetFromProc  chan BetInfo
	NewRound     chan RoundInfo
	WinnerName   chan WinInfo
	RoundResults chan ResultsInfo
}

// Function, which inits all the chans
func TunnelInit() (t Tunnel) {
	t.MoveToProc = make(chan MoveInfo)
	t.MoveFromProc = make(chan MoveInfo)
	t.BetToProc = make(chan BetInfo)
	t.BetFromProc = make(chan BetInfo)
	t.NewRound = make(chan RoundInfo)
	t.WinnerName = make(chan WinInfo)
	t.RoundResults = make(chan ResultsInfo)
	return
}

// Method, which closes all the chans, which provide info from "processor" to "handler"
func (t *Tunnel) CloseFromProcessor() {
	close(t.BetFromProc)
	close(t.MoveFromProc)
	close(t.WinnerName)
	close(t.NewRound)
	close(t.RoundResults)
}

// Method, which closes all the chans, which provide info from "handler" to "processor"
func (t *Tunnel) CloseFromHandler() {
	close(t.BetToProc)
	close(t.MoveToProc)
}

// Array which stores information about cards per player for usual rounds.
var UsualRulesArray = []int32{1, 1, 1, 2, 3, 4, 5,
	6, 7, 8, 9, 10, 11, 12,
	12, 12, 11, 9, 7, 4, 2,
	1, 1, 1}
