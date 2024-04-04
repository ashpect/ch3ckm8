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

func (b *Board) getColour(pos uint64) bool {
	return pos&b.whitePieces != 0
}

func (b *Board) piecePosToNotation(pos uint64) string {

	pieceType := b.getPieceType(pos)
	isWhite := b.getColour(pos)
	var base_string string = ""
	switch pieceType {
	case Pawn:
		base_string = ""
	case Knight:
		base_string = "N"
	case Bishop:
		base_string = "B"
	case Rook:
		base_string = "R"
	case Queen:
		base_string = "Q"
	case King:
		base_string = "K"
	}
	if !isWhite {
		base_string = strings.ToLower(base_string)
	}
	return base_string
}

func (b *Board) posToNotation(pos uint64) string {
	var posNotation string = ""
	var i uint64
	for i = 0x8000000000000000; i > 0; i >>= 1 {
		if pos&i != 0 {
			posNotation = string(rune(97+(63-i)%8)) + string(rune(49+int((63-i)/8)))
			break
		}
	}
	return posNotation
}

func (b *Board) notationToPos(notation string) uint64 {
	return 1 << uint64(int(notation[1]-49)*8+int(7-notation[0]-97))
}

func (b *Board) notationToPieceType(notation string) PieceType {
	switch strings.ToUpper(notation) {
	case "P":
		return Pawn
	case "N":
		return Knight
	case "B":
		return Bishop
	case "R":
		return Rook
	case "Q":
		return Queen
	case "K":
		return King
	}
	return NoPiece
}
