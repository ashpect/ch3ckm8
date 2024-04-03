package main

import (
	"strconv"
	"strings"
	"unicode"
)

const (
	rightEdge uint64    = 0x0101010101010101
	Pawn      PieceType = iota + 1
	Knight
	Bishop
	Rook
	Queen
	King
)

type PieceType int
type Board struct {
	whitePawns   uint64
	whiteKnights uint64
	whiteBishops uint64
	whiteRooks   uint64
	whiteQueens  uint64
	whiteKing    uint64
	whitePieces  uint64

	blackPawns   uint64
	blackKnights uint64
	blackBishops uint64
	blackRooks   uint64
	blackQueens  uint64
	blackKing    uint64
	blackPieces  uint64

	allPieces uint64
}

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

func (b *Board) Print(bitBoard uint64) {
	var i uint64
	for i = 0x8000000000000000; i > 0; i >>= 1 {
		if b.whitePawns&i != 0 {
			print("P ")
		} else if b.whiteKnights&i != 0 {
			print("N ")
		} else if b.whiteBishops&i != 0 {
			print("B ")
		} else if b.whiteRooks&i != 0 {
			print("R ")
		} else if b.whiteQueens&i != 0 {
			print("Q ")
		} else if b.whiteKing&i != 0 {
			print("K ")
		} else if b.blackPawns&i != 0 {
			print("p ")
		} else if b.blackKnights&i != 0 {
			print("n ")
		} else if b.blackBishops&i != 0 {
			print("b ")
		} else if b.blackRooks&i != 0 {
			print("r ")
		} else if b.blackQueens&i != 0 {
			print("q ")
		} else if b.blackKing&i != 0 {
			print("k ")
		} else {
			if bitBoard&i != 0 {
				print("1 ")
			} else {
				print(". ")
			}
		}
		if i&rightEdge != 0 {
			println()
		}
	}
	println()
}

func (b *Board) movePiece(initPos, finalPos uint64, pieceType PieceType, isWhite bool) {

	if isWhite {
		b.whitePieces &= ^initPos
		b.whitePieces |= finalPos
		b.allPieces = b.whitePieces | b.blackPieces
		if pieceType == Pawn {

			b.whitePawns &= ^initPos
			b.whitePawns |= finalPos

		} else if pieceType == Knight {

			b.whiteKnights &= ^initPos
			b.whiteKnights |= finalPos

		} else if pieceType == Bishop {

			b.whiteBishops &= ^initPos
			b.whiteBishops |= finalPos

		} else if pieceType == Rook {
			b.whiteRooks &= ^initPos
			b.whiteRooks |= finalPos

		} else if pieceType == Queen {

			b.whiteQueens &= ^initPos
			b.whiteQueens |= finalPos

		} else if pieceType == King {

			b.whiteKing &= ^initPos
			b.whiteKing |= finalPos
		}

	} else {
		b.blackPieces &= ^initPos
		b.blackPieces |= finalPos
		b.allPieces = b.whitePieces | b.blackPieces
		if pieceType == Pawn {

			b.blackPawns &= ^initPos
			b.blackPawns |= finalPos

		} else if pieceType == Knight {

			b.blackKnights &= ^initPos
			b.blackKnights |= finalPos

		} else if pieceType == Bishop {

			b.blackBishops &= ^initPos
			b.blackBishops |= finalPos

		} else if pieceType == Rook {
			b.blackRooks &= ^initPos
			b.blackRooks |= finalPos

		} else if pieceType == Queen {

			b.blackQueens &= ^initPos
			b.blackQueens |= finalPos

		} else if pieceType == King {

			b.blackKing &= ^initPos
			b.blackKing |= finalPos
		}
	}
}

var isWhite bool = false

func main() {
	var b Board
	b.Empty()
	fen := "2k1r1r1/p2q1p2/B1p5/3b3p/3Q4/5PP1/PP3K1P/R2R4 b - - 1 22"

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
	b.Print(0)
}
