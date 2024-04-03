package engine

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	rightEdge     uint64 = 0x0101010101010101
	leftEdge      uint64 = 0x8080808080808080
	bottomEdge    uint64 = 0x00000000000000FF
	topEdge       uint64 = 0xFF00000000000000
	whitePawnInit uint64 = 0x000000000000FF00
	blackPawnInit uint64 = 0x00FF000000000000

	diagBackRightRay    uint8 = 1 << 0
	backRay             uint8 = 1 << 1
	diagBackLeftRay     uint8 = 1 << 2
	rightRay            uint8 = 1 << 3
	leftRay             uint8 = 1 << 4
	diagForwardRightRay uint8 = 1 << 5
	forwardRay          uint8 = 1 << 6
	diagForwardLeftRay  uint8 = 1 << 7

	diagRays     uint8 = diagBackRightRay | diagBackLeftRay | diagForwardRightRay | diagForwardLeftRay
	straightRays uint8 = backRay | rightRay | leftRay | forwardRay
	allRays      uint8 = diagRays | straightRays
)

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

func (b *Board) Init() {
	b.whitePawns = 0x000000000000FF00
	b.whiteKnights = 0x0000000000000042
	b.whiteBishops = 0x0000000000000024
	b.whiteRooks = 0x0000000000000081
	b.whiteQueens = 0x0000000000000008
	b.whiteKing = 0x0000000000000010
	b.whitePieces = 0x000000000000FFFF

	b.blackPawns = 0x00FF000000000000
	b.blackKnights = 0x4200000000000000
	b.blackBishops = 0x2400000000000000
	b.blackRooks = 0x8100000000000000
	b.blackQueens = 0x0800000000000000
	b.blackKing = 0x1000000000000000
	b.blackPieces = 0xFFFF0000000000

	b.allPieces = 0xFFFF00000000FFFF
}

func printBitBoard(bitBoard uint64) {
	var i uint64
	for i = 0x8000000000000000; i > 0; i >>= 1 {
		if bitBoard&i != 0 {
			fmt.Print("1")
		} else {
			fmt.Print("0")
		}
		if i&0x0101010101010101 != 0 {
			fmt.Println()
		}
	}
}

func getMoves(piece uint64, depth int, rays uint8, b *Board, isWhite bool) uint64 {

	var moves, pre_comp, sameColorPieces uint64

	if isWhite {
		sameColorPieces = b.whitePieces
	} else {
		sameColorPieces = b.blackPieces
	}

	var rightCheck uint64 = rightEdge | b.allPieces
	var leftCheck uint64 = leftEdge | b.allPieces
	var topCheck uint64 = topEdge | b.allPieces
	var bottomCheck uint64 = bottomEdge | b.allPieces

	// Exploring all given directions symmetrically
	// and stopping exploring a direction when a piece is found or an edge is hit
	// or the depth is reached

	for i := 1; i < depth+1; i++ {
		if rays&diagBackRightRay != 0 {
			pre_comp = (piece & ^rightEdge >> (9 * i))
			if pre_comp&(rightCheck) != 0 || pre_comp == 0 {
				rays &= ^diagBackRightRay
			}
			moves |= pre_comp & ^sameColorPieces
		}
		if rays&backRay != 0 {
			pre_comp = (piece >> (8 * i))
			if pre_comp&(bottomCheck) != 0 || pre_comp == 0 {
				rays &= ^backRay
			}
			moves |= pre_comp & ^sameColorPieces
		}

		if rays&diagBackLeftRay != 0 {
			pre_comp = ((piece & ^leftEdge) >> (7 * i))
			if pre_comp&(leftCheck) != 0 || pre_comp == 0 {
				rays &= ^diagBackLeftRay
			}
			moves |= pre_comp & ^sameColorPieces
		}

		if rays&rightRay != 0 {
			pre_comp = ((piece &^ rightEdge) >> i)
			if pre_comp&(rightCheck) != 0 || pre_comp == 0 {
				rays &= ^rightRay
			}
			moves |= pre_comp & ^sameColorPieces
		}

		if rays&leftRay != 0 {
			pre_comp = ((piece & ^leftEdge) << i)
			if pre_comp&(leftCheck) != 0 || pre_comp == 0 {
				rays &= ^leftRay
			}
			moves |= pre_comp & ^sameColorPieces
		}

		if rays&diagForwardRightRay != 0 {
			pre_comp = ((piece & ^rightEdge) << (7 * i))
			if pre_comp&(rightCheck) != 0 || pre_comp == 0 {
				rays &= ^diagForwardRightRay
			}
			moves |= pre_comp & ^sameColorPieces
		}
		if rays&forwardRay != 0 {
			pre_comp = (piece << (8 * i))
			if pre_comp&(topCheck) != 0 || pre_comp == 0 {
				rays &= ^forwardRay
			}
			moves |= pre_comp & ^sameColorPieces
		}

		if rays&diagForwardLeftRay != 0 {
			pre_comp = ((piece & ^leftEdge) << (9 * i))
			if pre_comp&(leftCheck) != 0 || pre_comp == 0 {
				rays &= ^diagForwardLeftRay
			}
			moves |= pre_comp & ^sameColorPieces
		}
	}
	return moves
}

func getPawnMoves(piece uint64, b *Board, isWhite bool) uint64 {
	if isWhite {
		var moves uint64 = getMoves(piece, 1, diagForwardLeftRay|diagForwardRightRay, b, isWhite)
		if piece&whitePawnInit == 0 {
			moves |= getMoves(piece, 2, forwardRay, b, isWhite) &^ b.allPieces
		} else {
			moves |= getMoves(piece, 1, forwardRay, b, isWhite) &^ b.allPieces
		}
		return moves
	} else {
		var moves uint64 = getMoves(piece, 1, diagBackLeftRay|diagBackRightRay, b, isWhite)
		if piece&blackPawnInit == 0 {
			moves |= getMoves(piece, 2, backRay, b, isWhite) &^ b.allPieces
		} else {
			moves |= getMoves(piece, 1, backRay, b, isWhite) &^ b.allPieces
		}
		return moves
	}
	// TODO: Implement en passant
}

func getRookMoves(piece uint64, b *Board, isWhite bool) uint64 {
	var moves uint64 = getMoves(piece, 8, straightRays, b, isWhite)
	return moves
	// TODO: Implement castling
}
func getBishopMoves(piece uint64, b *Board, isWhite bool) uint64 {
	var moves uint64 = getMoves(piece, 8, diagRays, b, isWhite)
	return moves
}

func getQueenMoves(piece uint64, b *Board, isWhite bool) uint64 {
	var moves uint64 = getMoves(piece, 8, allRays, b, isWhite)
	return moves
}

func getKingMoves(piece uint64, b *Board, isWhite bool) uint64 {
	var moves uint64 = getMoves(piece, 1, allRays, b, isWhite)
	return moves
}
func evalPosition(b *Board) int {
	return 0
}
func getBestMove(b *Board) uint64 {
	return 0
}
func applyMove(b *Board, initPos, finalPos uint64, isWhite bool) {

}
func main() {
	var b Board
	var n_pieces int = 32

	for i := 0; i < n_pieces; i++ {
		b.allPieces |= (1 << rand.Intn(64))
	}

	printBitBoard(b.allPieces)
	start := time.Now()
	for i := 0; i < 100_000_000; i++ {
		getMoves(1<<rand.Intn(64), rand.Intn(7)+1, 0xFF, &b, true)
	}

	elapsed := time.Since(start)
	fmt.Printf("Time elapsed: %s\n", elapsed)
}
