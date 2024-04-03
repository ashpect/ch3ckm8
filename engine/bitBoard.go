package engine

import (
	"fmt"
)

// PieceType represents the type of a chess piece.
type PieceType int

const (
	bottomEdge       uint64 = 0x00000000000000FF
	bottomButOneEdge uint64 = 0x000000000000FF00

	rightEdge       uint64 = 0x0101010101010101
	rightButOneEdge uint64 = 0x0202020202020202

	leftEdge       uint64 = 0x8080808080808080
	leftButOneEdge uint64 = 0x4040404040404040

	topEdge       uint64 = 0xFF00000000000000
	topButOneEdge uint64 = 0x00FF000000000000

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

	Pawn PieceType = iota + 1
	Knight
	Bishop
	Rook
	Queen
	King
)

// Board represents the state of the chess board.
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

// printBitBoard prints the binary representation of a given bitboard.
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

// TODO: Implement threading
// getMoves calculates and returns the possible moves for a given chess piece on the board.
// It takes the piece's bitboard representation, the depth of exploration, the directions of possible movement,
// a pointer to the board, and a flag indicating whether the piece is white or not.
// It returns a bitboard representing the possible moves for the piece,
// where a 1 represents a possible move and a 0 represents an impossible move.
func getMoves(piece uint64, depth int, rays uint8, b *Board, isWhite bool) uint64 {

	var moves, pre_comp, sameColorPieces uint64

	if isWhite {
		sameColorPieces = b.whitePieces
	} else {
		sameColorPieces = b.blackPieces
	}

	// Defining the squares after reaching which, the piece will stop exploring the given direction
	// We don't need to check for top and bottom edges separately as bit-shifting will take care of it
	var rightCheck uint64 = rightEdge | b.allPieces
	var leftCheck uint64 = leftEdge | b.allPieces
	var topCheck uint64 = topEdge | b.allPieces
	var bottomCheck uint64 = bottomEdge | b.allPieces

	// Exploring all given directions symmetrically
	// and stopping exploring a direction when a piece is found or an edge is hit
	// or the depth is reached
	for i := 1; i < depth+1; i++ {
		if rays&diagBackRightRay != 0 {

			// Move the piece to the right and down unless it is on the right edge already
			pre_comp = (piece & ^rightEdge >> (9 * i))

			// If a piece is found or an edge is hit, stop exploring the direction
			if pre_comp&(rightCheck) != 0 || pre_comp == 0 {
				rays &= ^diagBackRightRay
			}

			// Add the possible moves to the moves bitboard
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

// getPawnMoves calculates and returns the possible moves for a pawn on the chessboard.
// It takes the pawn's position, the current board state, and a boolean flag indicating whether the pawn is white or not.
// If the pawn is white, it considers the forward and diagonal moves in the positive direction.
// If the pawn is black, it considers the forward and diagonal moves in the negative direction.
// The function also handles the special case of pawn's initial double move.
// It returns a bitboard representing the possible moves for the pawn.
// TODO: Implement en passant
func getPawnMoves(piece uint64, b *Board, isWhite bool) uint64 {
	if isWhite {
		var moves uint64 = getMoves(piece, 1, diagForwardLeftRay|diagForwardRightRay, b, isWhite)

		// If the pawn is on the second rank, it can move two squares forward
		if piece&bottomButOneEdge != 0 {
			moves |= getMoves(piece, 2, forwardRay, b, isWhite) &^ b.allPieces
		} else {
			moves |= getMoves(piece, 1, forwardRay, b, isWhite) &^ b.allPieces
		}
		return moves

	} else {
		var moves uint64 = getMoves(piece, 1, diagBackLeftRay|diagBackRightRay, b, isWhite)

		// If the pawn is on the seventh rank, it can move two squares forward
		if piece&topButOneEdge != 0 {
			moves |= getMoves(piece, 2, backRay, b, isWhite) &^ b.allPieces
		} else {
			moves |= getMoves(piece, 1, backRay, b, isWhite) &^ b.allPieces
		}
		return moves
	}
}

// TODO: Implement castling
// getRookMoves returns the possible moves for a rook piece on the given board.
// It takes the piece position, the board, and a flag indicating whether the piece is white or not.
// The function uses the getMoves helper function to calculate the possible moves in all four directions (up, down, left, right).
// The resulting moves are returned as a bitboard.
func getRookMoves(piece uint64, b *Board, isWhite bool) uint64 {
	var moves uint64 = getMoves(piece, 8, straightRays, b, isWhite)
	return moves
}

// getBishopMoves returns the possible moves for a bishop piece on the given board.
// It takes the position of the bishop piece, the board, and a flag indicating whether the bishop is white or not.
// The function calculates the possible moves by calling the getMoves function with the appropriate parameters.
// It returns a bitboard representing the possible moves for the bishop.
func getBishopMoves(piece uint64, b *Board, isWhite bool) uint64 {
	var moves uint64 = getMoves(piece, 8, diagRays, b, isWhite)
	return moves
}

// getQueenMoves returns the possible moves for a queen piece on the board.
// It takes the current position of the queen piece, the board state, and a flag indicating whether the piece is white or not.
// The function calculates and returns a bitboard representing the possible moves for the queen.
func getQueenMoves(piece uint64, b *Board, isWhite bool) uint64 {
	var moves uint64 = getMoves(piece, 8, allRays, b, isWhite)
	return moves
}

// getKnightMoves calculates the possible knight moves for a given piece on the board.
// It takes the piece's position, the board, and a flag indicating whether the piece is white or not.
// It returns a bitboard representing the possible knight moves.
func getKnightMoves(piece uint64, b *Board, isWhite bool) uint64 {
	var moves uint64
	var sameColorPieces uint64
	if isWhite {
		sameColorPieces = b.whitePieces
	} else {
		sameColorPieces = b.blackPieces
	}

	moves |= ((piece & ^rightEdge & ^rightButOneEdge) << 6) & ^sameColorPieces
	moves |= ((piece & ^rightEdge) << 15) &^ sameColorPieces
	moves |= ((piece & ^leftEdge) << 17) & ^sameColorPieces
	moves |= ((piece & ^leftEdge & ^leftButOneEdge) << 10) & ^sameColorPieces

	moves |= ((piece & ^leftEdge & ^leftButOneEdge) >> 6) & ^sameColorPieces
	moves |= ((piece & ^leftEdge) >> 15) & ^sameColorPieces
	moves |= ((piece & ^rightEdge) >> 17) & ^sameColorPieces
	moves |= ((piece & ^rightEdge & ^rightButOneEdge) >> 10) & ^sameColorPieces
	return moves
}

// getKingMoves returns the possible moves for a king piece on the given board.
// It takes the position of the king piece, the board, and a flag indicating whether the king is white or not.
// The function calculates the possible moves by using the getMoves function with the appropriate parameters.
// It returns a bitboard representing the possible moves for the king.
func getKingMoves(piece uint64, b *Board, isWhite bool) uint64 {
	var moves uint64 = getMoves(piece, 1, allRays, b, isWhite)
	return moves
}

func evalBoard(b *Board, isWhite bool) int {
	return 0
}

func getBestMove(b *Board) uint64 {
	return 0
}

// getPieceType returns the type of a given piece on the chess board.
// It takes a piece and a pointer to the Board struct as input.
// It checks the piece against the bitboards of each piece type and returns the corresponding PieceType.
// If the piece does not match any of the bitboards, it returns 0.
func getPieceType(piece uint64, b *Board) PieceType {
	if piece&b.whitePawns != 0 || piece&b.blackPawns != 0 {
		return Pawn
	} else if piece&b.whiteKnights != 0 || piece&b.blackKnights != 0 {
		return Knight
	} else if piece&b.whiteBishops != 0 || piece&b.blackBishops != 0 {
		return Bishop
	} else if piece&b.whiteRooks != 0 || piece&b.blackRooks != 0 {
		return Rook
	} else if piece&b.whiteQueens != 0 || piece&b.blackQueens != 0 {
		return Queen
	} else if piece&b.whiteKing != 0 || piece&b.blackKing != 0 {
		return King
	}
	return 0
}

// updatePieceTypeOnBoard updates the position of a piece on the board based on its initial and final positions.
// It takes a pointer to a Board struct, the initial position, final position, piece type, and a boolean indicating whether the piece is white.
// If the piece is white, it updates the corresponding white piece bitboard based on the piece type.
// If the piece is black, it updates the corresponding black piece bitboard based on the piece type.
func updatePieceTypeOnBoard(b *Board, initPos, finalPos uint64, pieceType PieceType, isWhite bool) {

	if isWhite {
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

// applyMove applies a move to the chess board.
// It updates the positions of the pieces on the board based on the initial and final positions provided.
// If the moving piece is white, it updates the whitePieces bitboard accordingly.
// If the moving piece is black, it updates the blackPieces bitboard accordingly.
// If a piece of the opposite color is taken during the move, it updates the respective bitboard and piece type.
// Finally, it updates the allPieces bitboard to reflect the new positions of all the pieces on the board.
func applyMove(b *Board, initPos, finalPos uint64, isWhite bool, pieceType PieceType) {
	if isWhite {
		b.whitePieces &= ^initPos
		b.whitePieces |= finalPos

		// If a black piece is taken
		if finalPos&b.blackPieces != 0 {
			b.blackPieces &= ^finalPos
			blackPieceType := getPieceType(finalPos, b)
			updatePieceTypeOnBoard(b, initPos, finalPos, blackPieceType, !isWhite)
		}
		updatePieceTypeOnBoard(b, initPos, finalPos, pieceType, isWhite)

	} else {
		b.blackPieces &= ^initPos
		b.blackPieces |= finalPos

		// If a white piece is taken
		if finalPos&b.whitePieces != 0 {
			b.whitePieces &= ^finalPos
			whitePieceType := getPieceType(finalPos, b)
			updatePieceTypeOnBoard(b, initPos, finalPos, whitePieceType, !isWhite)
		}
		updatePieceTypeOnBoard(b, initPos, finalPos, pieceType, isWhite)

	}
	b.allPieces = b.whitePieces | b.blackPieces
}

func isMoveValid(b *Board, initPos, finalPos uint64, isWhite bool) bool {
	return true
}

func main() {
	var b Board
	b.Init()
	printBitBoard(b.allPieces)
	printBitBoard(getKnightMoves(0x000000000200000, &b, true))

}
