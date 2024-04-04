package engine

import (
	"fmt"
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

func (b *Board) Print(isWhite bool) {
	if isWhite {
		var i uint64
		for i = 0x8000000000000000; i > 0; i >>= 1 {
			if b.whitePawns&i != 0 {
				fmt.Printf("P ")
			} else if b.whiteKnights&i != 0 {
				fmt.Printf("N ")
			} else if b.whiteBishops&i != 0 {
				fmt.Printf("B ")
			} else if b.whiteRooks&i != 0 {
				fmt.Printf("R ")
			} else if b.whiteQueens&i != 0 {
				fmt.Printf("Q ")
			} else if b.whiteKing&i != 0 {
				fmt.Printf("K ")
			} else if b.blackPawns&i != 0 {
				fmt.Printf("p ")
			} else if b.blackKnights&i != 0 {
				fmt.Printf("n ")
			} else if b.blackBishops&i != 0 {
				fmt.Printf("b ")
			} else if b.blackRooks&i != 0 {
				fmt.Printf("r ")
			} else if b.blackQueens&i != 0 {
				fmt.Printf("q ")
			} else if b.blackKing&i != 0 {
				fmt.Printf("k ")
			} else {
				fmt.Printf(". ")
			}
			if i&rightEdge != 0 {
				fmt.Println()
			}
		}
		fmt.Println()
	} else {
		var i uint64
		for i = 1; i <= 0x8000000000000000; i <<= 1 {
			if b.whitePawns&i != 0 {
				fmt.Printf("P ")
			} else if b.whiteKnights&i != 0 {
				fmt.Printf("N ")
			} else if b.whiteBishops&i != 0 {
				fmt.Printf("B ")
			} else if b.whiteRooks&i != 0 {
				fmt.Printf("R ")
			} else if b.whiteQueens&i != 0 {
				fmt.Printf("Q ")
			} else if b.whiteKing&i != 0 {
				fmt.Printf("K ")
			} else if b.blackPawns&i != 0 {
				fmt.Printf("p ")
			} else if b.blackKnights&i != 0 {
				fmt.Printf("n ")
			} else if b.blackBishops&i != 0 {
				fmt.Printf("b ")
			} else if b.blackRooks&i != 0 {
				fmt.Printf("r ")
			} else if b.blackQueens&i != 0 {
				fmt.Printf("q ")
			} else if b.blackKing&i != 0 {
				fmt.Printf("k ")
			} else {
				fmt.Printf(". ")
			}
			if i&leftEdge != 0 {
				fmt.Println()
			}
			if i == 0x8000000000000000 {
				break
			}
		}
		fmt.Println()
	}
}

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

// getPieceType returns the type of a given piece on the chess board.
// It takes a piece and a pointer to the Board struct as input.
// It checks the piece against the bitboards of each piece type and returns the corresponding PieceType.
// If the piece does not match any of the bitboards, it returns 0.
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
	return 0
}

func (b *Board) isCheck(isWhite bool) bool {
	oppMoves := b.getAllMoves(!isWhite)
	var kingPos uint64
	if isWhite {
		kingPos = b.whiteKing
	} else {
		kingPos = b.blackKing
	}
	for _, move := range oppMoves {
		for i := 0; i < 64; i++ {
			if move[1]&kingPos != 0 {
				return true
			}
		}
	}
	return false
}

func test() {
	var b Board
	b.Initialize()
	b.Print(true)
	moves := b.getAllLegalMoves(true)
	n_moves := 0
	for _, move := range moves {
		for i := 0; i < 64; i++ {
			if move[1]&(1<<uint64(i)) != 0 {
				wasPieceCaptured, capturePieceType := b.makeMove(move[0], 1<<uint64(i), true, b.getPieceType(move[0]))
				if wasPieceCaptured {
					fmt.Println("Captured: ", capturePieceType)
				}
				b.Print(true)
				fmt.Println(b.eval())
				n_moves++
				b.unmakeMove(move[0], 1<<uint64(i), true, wasPieceCaptured, capturePieceType)
			}
		}
	}
	fmt.Println(n_moves)
}

// check if the game has end or not
func (b *Board) isCheckmate() bool {
	return false
}
