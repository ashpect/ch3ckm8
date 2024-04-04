package engine

import (
	"strconv"
	"strings"
	"unicode"
)

// Init initializes the chess board with the starting positions of the pieces.
func (b *Board) Empty() {
	b.whitePawns = 0x0000000000000000
	b.whiteKnights = 0x0000000000000000
	b.whiteBishops = 0x0000000000000000
	b.whiteRooks = 0x0000000000000000
	b.whiteQueens = 0x0000000000000000
	b.whiteKing = 0x0000000000000000
	b.whitePieces = 0x0000000000000000

	b.blackPawns = 0x0000000000000000
	b.blackKnights = 0x0000000000000000
	b.blackBishops = 0x0000000000000000
	b.blackRooks = 0x0000000000000000
	b.blackQueens = 0x0000000000000000
	b.blackKing = 0x0000000000000000
	b.blackPieces = 0x0000000000000000

	b.allPieces = 0x0000000000000000
}

var isWhite bool = false

func parse(fen string) Board {

	var b Board
	b.Empty()

	// test fen = "2k1r1r1/p2q1p2/B1p5/3b3p/3Q4/5PP1/PP3K1P/R2R4 b - - 1 22"

	a := strings.Split(fen, " ")
	var curpos int = 63
	for i := 0; i < len(a[0]); i++ {
		chr := a[0][i]
		var pos uint64 = 1 << curpos
		if unicode.IsDigit(rune(chr)) || unicode.IsLetter(rune(chr)) {
			if chr >= 65 && chr <= 90 {
				isWhite = true
			}
			if unicode.ToLower(rune(chr)) == 'p' {
				b.movePiece(0, pos, Pawn, isWhite)
			} else if unicode.ToLower(rune(chr)) == 'r' {
				b.movePiece(0, pos, Rook, isWhite)
			} else if unicode.ToLower(rune(chr)) == 'n' {
				b.movePiece(0, pos, Knight, isWhite)
			} else if unicode.ToLower(rune(chr)) == 'b' {
				b.movePiece(0, pos, Bishop, isWhite)
			} else if unicode.ToLower(rune(chr)) == 'q' {
				b.movePiece(0, pos, Queen, isWhite)
			} else if unicode.ToLower(rune(chr)) == 'k' {
				b.movePiece(0, pos, King, isWhite)
			} else if chr >= 48 && chr <= 57 {
				num, _ := strconv.Atoi(string(chr))
				num2 := int(num)
				curpos -= num2 - 1
			}
			curpos--
		}
	}
	return b
}
