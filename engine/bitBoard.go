package engine

import (
	"fmt"
	"strings"
)

const (
	vertLine1 string = "    a   b   c   d   e   f   g   h  "
	vertLine2 string = "    h   g   f   e   d   c   b   a  "
	firstLine string = "  ┌───┬───┬───┬───┬───┬───┬───┬───┐"

	lineDelim string = "  ├───┼───┼───┼───┼───┼───┼───┼───┤"

	lastLine    string = "  └───┴───┴───┴───┴───┴───┴───┴───┘"
	colorNone   string = "\033[0m"
	colorYellow string = "\033[1;33m"
)

// Iinitializes the chess board with the starting positions of the pieces.
func (b *Board) Initialize() {
	b.whitePawns = 0x000000000000FF00
	b.whiteKnights = 0x0000000000000042
	b.whiteBishops = 0x0000000000000024
	b.whiteRooks = 0x0000000000000081
	b.whiteQueens = 0x0000000000000010
	b.whiteKing = 0x0000000000000008
	b.whitePieces = 0x000000000000FFFF

	b.blackPawns = 0x00FF000000000000
	b.blackKnights = 0x4200000000000000
	b.blackBishops = 0x2400000000000000
	b.blackRooks = 0x8100000000000000
	b.blackQueens = 0x1000000000000000
	b.blackKing = 0x0800000000000000
	b.blackPieces = 0xFFFF000000000000

	b.allPieces = 0xFFFF00000000FFFF
}

// Print prints the chess board to the console.
// It takes a boolean as input to determine if the board should be printed from the perspective of the white player.
func (b *Board) PrintBoard(isWhite bool, bestMov uint64) {
	if isWhite {
		var i uint64
		v := 8
		fmt.Println(firstLine)
		for i = 0x8000000000000000; i > 0; i >>= 1 {
			print(colorNone)
			if i&leftEdge != 0 {
				fmt.Print(v, " ")
				v--
			}
			fmt.Printf("│ ")
			if i&bestMov != 0 {
				print(colorYellow)
			}
			pieceTypeString := string(rune(b.getPieceType(i)))
			if b.whitePieces&i != 0 {
				fmt.Print(pieceTypeString)
			} else if b.blackPieces&i != 0 {
				fmt.Print(strings.ToLower(pieceTypeString))
			} else {
				fmt.Printf(" ")
			}
			fmt.Printf(" ")
			if i&rightEdge != 0 {
				fmt.Println("│")
				if i&bottomEdge != 0 {
					fmt.Println(lastLine)
					break
				} else {
					fmt.Println(lineDelim)
				}
			}
		}
		fmt.Println(vertLine1)
		fmt.Println()
	} else {
		var i uint64
		v := 1
		fmt.Println(firstLine)
		for i = 1; i <= 0x8000000000000000; i <<= 1 {
			print(colorNone)
			if i&rightEdge != 0 {
				fmt.Print(v, " ")
				v++
			}
			fmt.Printf("│ ")
			if i&bestMov != 0 {
				print(colorYellow)
			}
			pieceTypeString := string(rune(b.getPieceType(i)))
			if b.whitePieces&i != 0 {
				fmt.Print(pieceTypeString)
			} else if b.blackPieces&i != 0 {
				fmt.Print(strings.ToLower(pieceTypeString))
			} else {
				fmt.Printf(" ")
			}
			fmt.Printf(" ")
			if i&leftEdge != 0 {
				fmt.Println("│")
				if i&topEdge != 0 {
					fmt.Println(lastLine)
					break
				} else {
					fmt.Println(lineDelim)
				}
			}
			if i == 0x8000000000000000 {
				break
			}
		}
		fmt.Println(vertLine2)
		fmt.Println()
	}
}

// printBitBoard prints a given uint64 as a chess board.
func printBitBoard(bitBoard uint64) {
	var i uint64
	for i = 0x8000000000000000; i > 0; i >>= 1 {
		if bitBoard&i != 0 {
			fmt.Printf("1 ")
		} else {
			fmt.Printf(". ")
		}
		if i&rightEdge != 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}

// getColour returns the colour of the piece at a given position.
func (b *Board) getColour(pos uint64) bool {
	return pos&b.whitePieces != 0
}

// getPieceType returns the type of the piece at a given position.
func (b *Board) getPieceType(pos uint64) PieceType {
	if pos&b.whitePawns != 0 || pos&b.blackPawns != 0 {
		return Pawn
	} else if pos&b.whiteKnights != 0 || pos&b.blackKnights != 0 {
		return Knight
	} else if pos&b.whiteBishops != 0 || pos&b.blackBishops != 0 {
		return Bishop
	} else if pos&b.whiteRooks != 0 || pos&b.blackRooks != 0 {
		return Rook
	} else if pos&b.whiteQueens != 0 || pos&b.blackQueens != 0 {
		return Queen
	} else if pos&b.whiteKing != 0 || pos&b.blackKing != 0 {
		return King
	}
	return NoPiece
}

// check if a king is currently under attack
func (b *Board) isCheck(isWhite bool) bool {
	oppMoves := b.getAllMoves(!isWhite)
	var kingPos uint64
	if isWhite {
		kingPos = b.whiteKing
	} else {
		kingPos = b.blackKing
	}
	for _, move := range oppMoves {
		if move[1]&kingPos != 0 {
			return true
		}
	}
	return false
}

// check if the game has ended or not
func (b *Board) isCheckmate(isWhite bool) bool {
	if b.isCheck(isWhite) {
		if len(b.getAllLegalMoves(isWhite)) == 0 {
			return true
		}
	}
	return false
}

func test() {
	var b Board = parse("7k/3p3P/4npPK/2N5/8/8/8/8 w - - 0 1")
	b.makeMove(0x2000000000, 0x080000000000, true, Knight)
	moves := b.getAllLegalMoves(false)
	for _, move := range moves {
		for i := 0; i < 64; i++ {
			if move[1]&(1<<uint(i)) != 0 {
				printBitBoard(1 << uint(i))
			}
		}
	}
	b.PrintBoard(true, 0)
}
