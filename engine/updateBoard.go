package engine

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

// makeMove applies a move to the chess board.
// It updates the positions of the pieces on the board based on the initial and final positions provided.
// If the moving piece is white, it updates the whitePieces bitboard accordingly.
// If the moving piece is black, it updates the blackPieces bitboard accordingly.
// If a piece of the opposite color is taken during the move, it updates the respective bitboard and piece type.
// Finally, it updates the allPieces bitboard to reflect the new positions of all the pieces on the board.

// returns a boolean indicating whether a piece was captured during the move and the type of the captured piece.
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
