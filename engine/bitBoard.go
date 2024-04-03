package engine

import (
	"fmt"
	"time"
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

	diagBackRightDir    uint8 = 1 << 0
	backDir             uint8 = 1 << 1
	diagBackLeftDir     uint8 = 1 << 2
	rightDir            uint8 = 1 << 3
	leftDir             uint8 = 1 << 4
	diagForwardRightDir uint8 = 1 << 5
	forwardDir          uint8 = 1 << 6
	diagForwardLeftDir  uint8 = 1 << 7

	diagDirs     uint8 = diagBackRightDir | diagBackLeftDir | diagForwardRightDir | diagForwardLeftDir
	straightDirs uint8 = backDir | rightDir | leftDir | forwardDir
	allDirs      uint8 = diagDirs | straightDirs

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
	canWhiteCastle bool
	
	blackPawns   uint64
	blackKnights uint64
	blackBishops uint64
	blackRooks   uint64
	blackQueens  uint64
	blackKing    uint64
	blackPieces  uint64
	canBlackCastle bool

	allPieces uint64
}

// Init initializes the chess board with the starting positions of the pieces.
func (b *Board) Init() {
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
func printBitBoard(bitBoard uint64) {
	var i uint64
	for i = 0x8000000000000000; i > 0; i >>= 1 {
		if bitBoard&i != 0 {
			print("1 ")
		} else {
			print(". ")
		}
		if i&rightEdge != 0 {
			println()
		}
	}
	println()
}

// TODO: Implement threading
// getMoves calculates and returns the possible moves for a given chess piece on the board.
// It takes the piece's bitboard representation, the depth of exploration, the directions of possible movement,
// a pointer to the board, and a flag indicating whether the piece is white or not.
// It returns a bitboard representing the possible moves for the piece,
// where a 1 represents a possible move and a 0 represents an impossible move.
func (b *Board) getMoves(piece uint64, depth int, dirs uint8, isWhite bool) uint64 {

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
		if dirs&diagBackRightDir != 0 {

			// Move the piece to the right and down unless it is on the right edge already
			pre_comp = (piece & ^rightEdge >> (9 * i))

			// If a piece is found or an edge is hit, stop exploring the direction
			if pre_comp&(rightCheck) != 0 || pre_comp == 0 {
				dirs &= ^diagBackRightDir
			}

			// Add the possible moves to the moves bitboard
			moves |= pre_comp & ^sameColorPieces
		}

		if dirs&backDir != 0 {
			pre_comp = (piece >> (8 * i))
			if pre_comp&(bottomCheck) != 0 || pre_comp == 0 {
				dirs &= ^backDir
			}
			moves |= pre_comp & ^sameColorPieces
		}

		if dirs&diagBackLeftDir != 0 {
			pre_comp = ((piece & ^leftEdge) >> (7 * i))
			if pre_comp&(leftCheck) != 0 || pre_comp == 0 {
				dirs &= ^diagBackLeftDir
			}
			moves |= pre_comp & ^sameColorPieces
		}

		if dirs&rightDir != 0 {
			pre_comp = ((piece &^ rightEdge) >> i)
			if pre_comp&(rightCheck) != 0 || pre_comp == 0 {
				dirs &= ^rightDir
			}
			moves |= pre_comp & ^sameColorPieces
		}

		if dirs&leftDir != 0 {
			pre_comp = ((piece & ^leftEdge) << i)
			if pre_comp&(leftCheck) != 0 || pre_comp == 0 {
				dirs &= ^leftDir
			}
			moves |= pre_comp & ^sameColorPieces
		}

		if dirs&diagForwardRightDir != 0 {
			pre_comp = ((piece & ^rightEdge) << (7 * i))
			if pre_comp&(rightCheck) != 0 || pre_comp == 0 {
				dirs &= ^diagForwardRightDir
			}
			moves |= pre_comp & ^sameColorPieces
		}
		if dirs&forwardDir != 0 {
			pre_comp = (piece << (8 * i))
			if pre_comp&(topCheck) != 0 || pre_comp == 0 {
				dirs &= ^forwardDir
			}
			moves |= pre_comp & ^sameColorPieces
		}

		if dirs&diagForwardLeftDir != 0 {
			pre_comp = ((piece & ^leftEdge) << (9 * i))
			if pre_comp&(leftCheck) != 0 || pre_comp == 0 {
				dirs &= ^diagForwardLeftDir
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
func (b *Board) getPawnMoves(piece uint64, isWhite bool) uint64 {
	if isWhite {
		var moves uint64 = b.getMoves(piece, 1, diagForwardLeftDir|diagForwardRightDir, isWhite) & b.blackPieces

		// If the pawn is on the second rank, it can move two squares forward
		if piece&bottomButOneEdge != 0 {
			moves |= b.getMoves(piece, 2, forwardDir, isWhite) &^ b.allPieces
		} else {
			moves |= b.getMoves(piece, 1, forwardDir, isWhite) &^ b.allPieces
		}
		return moves

	} else {
		var moves uint64 = b.getMoves(piece, 1, diagBackLeftDir|diagBackRightDir, isWhite) & b.whitePieces

		// If the pawn is on the seventh rank, it can move two squares forward
		if piece&topButOneEdge != 0 {
			moves |= b.getMoves(piece, 2, backDir, isWhite) &^ b.allPieces
		} else {
			moves |= b.getMoves(piece, 1, backDir, isWhite) &^ b.allPieces
		}
		return moves
	}
}

// TODO: Implement castling
// getRookMoves returns the possible moves for a rook piece on the given board.
// It takes the piece position, the board, and a flag indicating whether the piece is white or not.
// The function uses the getMoves helper function to calculate the possible moves in all four directions (up, down, left, right).
// The resulting moves are returned as a bitboard.
func (b *Board) getRookMoves(piece uint64, isWhite bool) uint64 {
	var moves uint64 = b.getMoves(piece, 7, straightDirs, isWhite)
	return moves
}

// getBishopMoves returns the possible moves for a bishop piece on the given board.
// It takes the position of the bishop piece, the board, and a flag indicating whether the bishop is white or not.
// The function calculates the possible moves by calling the getMoves function with the appropriate parameters.
// It returns a bitboard representing the possible moves for the bishop.
func (b *Board) getBishopMoves(piece uint64, isWhite bool) uint64 {
	var moves uint64 = b.getMoves(piece, 7, diagDirs, isWhite)
	return moves
}

// getQueenMoves returns the possible moves for a queen piece on the board.
// It takes the current position of the queen piece, the board state, and a flag indicating whether the piece is white or not.
// The function calculates and returns a bitboard representing the possible moves for the queen.
func (b *Board) getQueenMoves(piece uint64, isWhite bool) uint64 {
	var moves uint64 = b.getMoves(piece, 7, allDirs, isWhite)
	return moves
}

// getKnightMoves calculates the possible knight moves for a given piece on the board.
// It takes the piece's position, the board, and a flag indicating whether the piece is white or not.
// It returns a bitboard representing the possible knight moves.
func (b *Board) getKnightMoves(piece uint64, isWhite bool) uint64 {
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
func (b *Board) getKingMoves(piece uint64, isWhite bool) uint64 {
	var moves uint64 = b.getMoves(piece, 1, allDirs, isWhite)
	return moves
}

func (b *Board) eval() int {
	return 0
}

func (b *Board) getBestMove(isWhite bool) uint64 {
	return 0
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

// updatePieceTypeOnBoard updates the position of a piece on the board based on its initial and final positions.
// It takes a pointer to a Board struct, the initial position, final position, piece type, and a boolean indicating whether the piece is white.
// If the piece is white, it updates the corresponding white piece bitboard based on the piece type.
// If the piece is black, it updates the corresponding black piece bitboard based on the piece type.
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
func (b *Board) getAllLegalMoves(isWhite bool) [][2]uint64 {
	var moves [][2]uint64
	if isWhite {
		for i := 0; i < 64; i++ {
			var cur_pos uint64 = 1 << uint64(i)
			if b.whitePawns&cur_pos != 0 {
				moves = append(moves, [2]uint64{cur_pos, b.getPawnMoves(cur_pos, isWhite)})
			} else if b.whiteKnights&cur_pos != 0 {
				moves = append(moves, [2]uint64{cur_pos, b.getKnightMoves(cur_pos, isWhite)})
			} else if b.whiteBishops&cur_pos != 0 {
				moves = append(moves, [2]uint64{cur_pos, b.getBishopMoves(cur_pos, isWhite)})
			} else if b.whiteRooks&cur_pos != 0 {
				moves = append(moves, [2]uint64{cur_pos, b.getRookMoves(cur_pos, isWhite)})
			} else if b.whiteQueens&cur_pos != 0 {
				moves = append(moves, [2]uint64{cur_pos, b.getQueenMoves(cur_pos, isWhite)})
			} else if b.whiteKing&cur_pos != 0 {
				moves = append(moves, [2]uint64{cur_pos, b.getKingMoves(cur_pos, isWhite)})
			}
		}
	} else {
		for i := 0; i < 64; i++ {
			var cur_pos uint64 = 1 << uint64(i)
			if b.blackPawns&cur_pos != 0 {
				moves = append(moves, [2]uint64{cur_pos, b.getPawnMoves(cur_pos, isWhite)})
			} else if b.blackKnights&uint64(cur_pos) != 0 {
				moves = append(moves, [2]uint64{cur_pos, b.getKnightMoves(cur_pos, isWhite)})
			} else if b.blackBishops&uint64(cur_pos) != 0 {
				moves = append(moves, [2]uint64{cur_pos, b.getBishopMoves(cur_pos, isWhite)})
			} else if b.blackRooks&uint64(cur_pos) != 0 {
				moves = append(moves, [2]uint64{cur_pos, b.getRookMoves(cur_pos, isWhite)})
			} else if b.blackQueens&uint64(cur_pos) != 0 {
				moves = append(moves, [2]uint64{cur_pos, b.getQueenMoves(cur_pos, isWhite)})
			} else if b.blackKing&uint64(cur_pos) != 0 {
				moves = append(moves, [2]uint64{cur_pos, b.getKingMoves(cur_pos, isWhite)})
			}
		}
	}
	return moves
}

// makeMove applies a move to the chess board.
// It updates the positions of the pieces on the board based on the initial and final positions provided.
// If the moving piece is white, it updates the whitePieces bitboard accordingly.
// If the moving piece is black, it updates the blackPieces bitboard accordingly.
// If a piece of the opposite color is taken during the move, it updates the respective bitboard and piece type.
// Finally, it updates the allPieces bitboard to reflect the new positions of all the pieces on the board.
func (b *Board) makeMove(initPos, finalPos uint64, isWhite bool, pieceType PieceType) (bool, PieceType) {
	if isWhite {
		// If a black piece is taken
		var wasPieceCaptured bool = (finalPos&b.blackPieces != 0)
		var capturedPieceType PieceType = 0
		if wasPieceCaptured {
			b.blackPieces &= ^finalPos
			capturedPieceType = b.getPieceType(finalPos)
			b.movePiece(finalPos, 0, capturedPieceType, !isWhite)
		}
		b.movePiece(initPos, finalPos, pieceType, isWhite)
		return wasPieceCaptured, capturedPieceType

	} else {
		// If a white piece is taken
		var wasPieceCaptured bool = (finalPos&b.whitePieces != 0)
		var capturedPieceType PieceType = 0
		if wasPieceCaptured {
			b.whitePieces &= ^finalPos
			capturedPieceType = b.getPieceType(finalPos)
			b.movePiece(finalPos, 0, capturedPieceType, !isWhite)
		}
		b.movePiece(initPos, finalPos, pieceType, isWhite)
		return wasPieceCaptured, capturedPieceType
	}
}

// unmakeMove takes the initial position, final position, and other parameters of a move and reverts the board state to the previous state.
// Parameters:
// - initPos: The initial position of the piece being moved.
// - finalPos: The final position of the piece being moved.
// - wasPieceCaptured: A boolean indicating whether a piece was captured during the move.
// - isWhite: A boolean indicating whether the moving piece is white.
// - capturedPieceType: The type of the captured piece, if any.
func (b *Board) unmakeMove(initPos, finalPos uint64, isWhite, wasPieceCaptured bool, capturedPieceType PieceType) {
	if isWhite {
		b.movePiece(finalPos, initPos, b.getPieceType(finalPos), isWhite)
		if wasPieceCaptured {
			b.movePiece(0, finalPos, capturedPieceType, !isWhite)
		}
	} else {
		b.movePiece(finalPos, initPos, b.getPieceType(finalPos), isWhite)
		if wasPieceCaptured {
			b.movePiece(0, finalPos, capturedPieceType, !isWhite)
		}
	}
	b.allPieces = b.whitePieces | b.blackPieces
}
func (b *Board) isMoveValid(initPos, finalPos uint64, isWhite bool) bool {
	return true
}

func test() {
	var b Board
	b.Init()
	b.makeMove(0x0000000000001000, 0, true, b.getPieceType(0x0000000000001000))
	moves := b.getAllLegalMoves(true)
	n_moves := 0
	for _, move := range moves {
		for i := 0; i < 64; i++ {
			if move[1]&(1<<uint64(i)) != 0 {
				wasPieceCaptured, capturePieceType := b.makeMove(move[0], 1<<uint64(i), true, b.getPieceType(move[0]))
				if wasPieceCaptured {
					println("Captured: ", capturePieceType)
				}
				n_moves++
				b.unmakeMove(move[0], 1<<uint64(i), true, wasPieceCaptured, capturePieceType)
			}
		}
	}
}
func main() {
	start := time.Now()
	for i := 0; i < 100_000_000; i++ {
		test()
	}
	elapsed := time.Since(start)
	fmt.Printf("Time: %s", elapsed)
}
