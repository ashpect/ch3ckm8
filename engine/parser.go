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

func parse(fen string) Board {

	var b Board
	b.Empty()
	var isWhite bool = false

	a := strings.Split(fen, " ")[0]
	var curpos int = 63
	for i := 0; i < len(a); i++ {
		chr := a[i]
		var pos uint64 = 1 << uint64(curpos)
		if unicode.IsDigit(rune(chr)) || unicode.IsLetter(rune(chr)) {
			if chr >= 65 && chr <= 90 {
				isWhite = true
			} else {
				isWhite = false
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

func (b *Board) piecePosToNotation(pos uint64) string {

	pieceType := b.getPieceType(pos)
	isWhite := b.getColour(pos)
	var base_string string = ""
	switch pieceType {
	case NoPiece:
		base_string = ""
	case Pawn:
		base_string = "P"
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
	for i := 63; i >= 0; i-- {
		var cur_pos uint64 = 1 << uint64(i)
		if pos&cur_pos != 0 {
			file := int((63 - i) % 8)
			rank := int((i) / 8)
			posNotation = string(rune(97+file)) + string(rune(49+rank))
			break
		}
	}
	return posNotation
}

func flipNotation(notation string) string {
	file := int(notation[0] - 97)
	rank := int(notation[1] - 49)
	return string(rune(97+7-file)) + string(rune(49+7-rank))
}

func (b *Board) notationToPos(notation string) uint64 {
	file := int(notation[0] - 97) // columns
	rank := int(notation[1] - 49) // rows
	pos := uint64(rank*8 + (7 - file))
	return 1 << pos
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
