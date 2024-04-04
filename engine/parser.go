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
			if chr >= 48 && chr <= 57 {
				num, _ := strconv.Atoi(string(chr))
				num2 := int(num)
				curpos -= num2 - 1
			} else {
				b.movePiece(0, pos, PieceType(unicode.ToUpper(rune(chr))), isWhite)
			}
			curpos--
		}
	}
	return b
}

func (b *Board) piecePosToNotation(pos uint64) string {

	return string(rune(b.getPieceType(pos)))

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

func (b *Board) notationToMove(notation string) (PieceType, uint64, uint64) {
	return PieceType(rune(notation[0])), b.notationToPos(notation[1:3]), b.notationToPos(notation[4:])
}

func (b *Board) moveToNotation(isWhite bool, initPos, finalPos uint64, pieceType PieceType, wasPieceCaptured bool) string {
	initPosNot := b.posToNotation(initPos)
	finalPosNot := b.posToNotation(finalPos)
	if !isWhite {
		initPosNot = flipNotation(initPosNot)
		finalPosNot = flipNotation(finalPosNot)
	}
	var sep string = "-"
	if wasPieceCaptured {
		sep = "x"
	}
	return string(rune(pieceType)) + initPosNot + sep + finalPosNot
}
