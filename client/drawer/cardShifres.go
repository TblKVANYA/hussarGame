// package drawer provides a console client for hussar.
// Functions are:
// func Init       (int32, int32)
// func UpdateIndex(int32)
// func SetBet     (int32, int32)
// func SetDecks   ([]int32)
// func SetCard    (int32, int32)
// func SetTrump   (int32)
// func SetArrow   (int32)
// func FlushCard  (int32, int32)
// func FlushBoard ()
// func FlushRound ()
// func AddWin     (int32)
// func AddRes     (int32, int32)
// func ReplaceCard(int32, int32)
// func Draw       ()
package drawer

var (
	ranks = map[int32]string{
		0:  "5",
		1:  "6",
		2:  "7",
		3:  "8",
		4:  "9",
		5:  "10",
		6:  "J",
		7:  "Q",
		8:  "K",
		9:  "A",
		10: "S", // superior
	}

	suits = map[int32]string{
		0: "♡",
		1: "♢",
		2: "♤",
		3: "♧",
		4: "☓", // empty suit
	}
)

// GetMagic returns the integer code of J♧
func GetMagic() int32 {
	return 27
}
