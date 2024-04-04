package engine

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

// TODO: Implement threading
// getMoves calculates and returns the possible moves for a given chess piece on the board.
// It takes the piece's bitboard representation, the depth of exploration, the directions of possible movement,
// a pointer to the board, and a flag indicating whether the piece is white or not.
// It returns a bitboard representing the possible moves for the piece,
// where a 1 represents a possible move and a 0 represents an impossible move.
func (b *Board) getMoves(position uint64, depth int, dirs uint8, isWhite bool) uint64 {

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
			pre_comp = (position & ^rightEdge >> (9 * i))

			// If a piece is found or an edge is hit, stop exploring the direction
			if pre_comp&(rightCheck) != 0 || pre_comp == 0 {
				dirs &= ^diagBackRightDir
			}

			// Add the possible moves to the moves bitboard
			moves |= pre_comp & ^sameColorPieces
		}

		if dirs&backDir != 0 {
			pre_comp = (position >> (8 * i))
			if pre_comp&(bottomCheck) != 0 || pre_comp == 0 {
				dirs &= ^backDir
			}
			moves |= pre_comp & ^sameColorPieces
		}

		if dirs&diagBackLeftDir != 0 {
			pre_comp = ((position & ^leftEdge) >> (7 * i))
			if pre_comp&(leftCheck) != 0 || pre_comp == 0 {
				dirs &= ^diagBackLeftDir
			}
			moves |= pre_comp & ^sameColorPieces
		}

		if dirs&rightDir != 0 {
			pre_comp = ((position &^ rightEdge) >> i)
			if pre_comp&(rightCheck) != 0 || pre_comp == 0 {
				dirs &= ^rightDir
			}
			moves |= pre_comp & ^sameColorPieces
		}

		if dirs&leftDir != 0 {
			pre_comp = ((position & ^leftEdge) << i)
			if pre_comp&(leftCheck) != 0 || pre_comp == 0 {
				dirs &= ^leftDir
			}
			moves |= pre_comp & ^sameColorPieces
		}

		if dirs&diagForwardRightDir != 0 {
			pre_comp = ((position & ^rightEdge) << (7 * i))
			if pre_comp&(rightCheck) != 0 || pre_comp == 0 {
				dirs &= ^diagForwardRightDir
			}
			moves |= pre_comp & ^sameColorPieces
		}
		if dirs&forwardDir != 0 {
			pre_comp = (position << (8 * i))
			if pre_comp&(topCheck) != 0 || pre_comp == 0 {
				dirs &= ^forwardDir
			}
			moves |= pre_comp & ^sameColorPieces
		}

		if dirs&diagForwardLeftDir != 0 {
			pre_comp = ((position & ^leftEdge) << (9 * i))
			if pre_comp&(leftCheck) != 0 || pre_comp == 0 {
				dirs &= ^diagForwardLeftDir
			}
			moves |= pre_comp & ^sameColorPieces
		}
	}
	return moves
}

// DO
func (b *Board) getBestMove(isWhite bool) uint64 {
	if isWhite {
		return 0
	} else {
		return 1
	}
}

func (b *Board) getAllMoves(isWhite bool) [][2]uint64 {
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
			} else if b.blackKnights&cur_pos != 0 {
				moves = append(moves, [2]uint64{cur_pos, b.getKnightMoves(cur_pos, isWhite)})
			} else if b.blackBishops&cur_pos != 0 {
				moves = append(moves, [2]uint64{cur_pos, b.getBishopMoves(cur_pos, isWhite)})
			} else if b.blackRooks&cur_pos != 0 {
				moves = append(moves, [2]uint64{cur_pos, b.getRookMoves(cur_pos, isWhite)})
			} else if b.blackQueens&cur_pos != 0 {
				moves = append(moves, [2]uint64{cur_pos, b.getQueenMoves(cur_pos, isWhite)})
			} else if b.blackKing&cur_pos != 0 {
				moves = append(moves, [2]uint64{cur_pos, b.getKingMoves(cur_pos, isWhite)})
			}
		}
	}
	return moves
}

func (b *Board) getAllLegalMoves(isWhite bool) [][2]uint64 {
	allMoves := b.getAllMoves(isWhite)
	for _, move := range allMoves {
		for i := 0; i < 64; i++ {
			var cur_pos uint64
			cur_pos = 1 << uint64(i)
			if move[1]&cur_pos != 0 {
				wasPieceCaptured, capturedPieceType := b.makeMove(move[0], cur_pos, isWhite, b.getPieceType(move[0]))
				if b.isCheck(isWhite) {
					move[1] &= ^cur_pos
				}
				b.unmakeMove(move[0], cur_pos, isWhite, wasPieceCaptured, capturedPieceType)
			}
		}
	}
	return allMoves
}
