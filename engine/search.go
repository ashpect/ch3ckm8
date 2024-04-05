package engine

func (b *Board) alphaBetaMiniMax(isWhite bool, alpha, beta float64, depth int) (float64, [2]uint64) {
	if depth == 0 {
		a, c := b.eval()
		return a + c, [2]uint64{0, 0}
	}
	if isWhite {
		lastMove := [2]uint64{0, 0}
		bestMove := [2]uint64{0, 0}
		for _, move := range b.getAllLegalMoves(isWhite) {
			var to_break bool = false
			for i := 0; i < 64; i++ {
				var cur_pos uint64 = 1 << uint64(i)
				if move[1]&cur_pos != 0 {
					lastMove = [2]uint64{move[0], cur_pos}

					//
					var wasPieceCaptured bool
					var capturedPieceType PieceType
					var kingmove [2]uint64 = [2]uint64{0, 0}
					isCastleMove, valueReceived := checkMoveIsCastle(move)

					if isCastleMove {
						b.makeMove(b.notationToPos(castleMoveInfo[valueReceived].KingMove[1:3]), cur_pos, isWhite, King)
						b.makeMove(b.notationToPos(castleMoveInfo[valueReceived].RookMove[3:5]), cur_pos, isWhite, Rook)
						castleMoveInfo[valueReceived].CanCastle = false
					} else {
						wasPieceCaptured, capturedPieceType = b.makeMove(move[0], cur_pos, isWhite, b.getPieceType(move[0]))
						if (kingmove != [2]uint64{0, 0}) && (b.getPieceType(move[0]) == King) {
							kingmove = move
							castleMoveInfo["O-O"].CanCastle = false
							castleMoveInfo["O-O-O"].CanCastle = false
						}
					}

					moveAlpha, _ := b.alphaBetaMiniMax(!isWhite, alpha, beta, depth-1)
					if isCastleMove {
						b.unmakeMove(b.notationToPos(castleMoveInfo[valueReceived].KingMove[1:3]), cur_pos, isWhite, false, NoPiece)
						b.unmakeMove(b.notationToPos(castleMoveInfo[valueReceived].RookMove[3:5]), cur_pos, isWhite, false, NoPiece)
						castleMoveInfo[valueReceived].CanCastle = true
					} else {
						b.unmakeMove(move[0], cur_pos, isWhite, wasPieceCaptured, capturedPieceType)
						if kingmove == move {
							kingmove = [2]uint64{0, 0}
							castleMoveInfo["O-O"].CanCastle = true
							castleMoveInfo["O-O-O"].CanCastle = true
						}
					}
					//

					if moveAlpha > alpha {
						alpha = moveAlpha
						bestMove = [2]uint64{move[0], cur_pos}
					}
					if beta <= alpha {
						to_break = true
						break
					}
				}
			}
			if to_break {
				break
			}
		}
		if bestMove == [2]uint64{0, 0} {
			bestMove = lastMove
		}

		return alpha, bestMove
	} else {
		lastMove := [2]uint64{0, 0}
		bestMove := [2]uint64{0, 0}
		for _, move := range b.getAllLegalMoves(isWhite) {
			var to_break bool = false
			for i := 0; i < 64; i++ {
				var cur_pos uint64 = 1 << uint64(i)
				if move[1]&cur_pos != 0 {
					lastMove = [2]uint64{move[0], cur_pos}

					//
					var wasPieceCaptured bool
					var capturedPieceType PieceType
					var kingmove [2]uint64 = [2]uint64{0, 0}
					isCastleMove, valueReceived := checkMoveIsCastle(move)
					if isCastleMove {
						b.makeMove(b.notationToPos(castleMoveInfo[valueReceived].KingMove[1:3]), cur_pos, isWhite, King)
						b.makeMove(b.notationToPos(castleMoveInfo[valueReceived].RookMove[3:5]), cur_pos, isWhite, Rook)
						castleMoveInfo[valueReceived].CanCastle = false
					} else {
						wasPieceCaptured, capturedPieceType = b.makeMove(move[0], cur_pos, isWhite, b.getPieceType(move[0]))
						if (kingmove != [2]uint64{0, 0}) && (b.getPieceType(move[0]) == King) {
							kingmove = move
							castleMoveInfo["o-o"].CanCastle = false
							castleMoveInfo["o-o-o"].CanCastle = false
						}
					}
					moveBeta, _ := b.alphaBetaMiniMax(!isWhite, alpha, beta, depth-1)
					if isCastleMove {
						b.unmakeMove(b.notationToPos(castleMoveInfo[valueReceived].KingMove[1:3]), cur_pos, isWhite, false, NoPiece)
						b.unmakeMove(b.notationToPos(castleMoveInfo[valueReceived].RookMove[3:5]), cur_pos, isWhite, false, NoPiece)
						castleMoveInfo[valueReceived].CanCastle = true
					} else {
						b.unmakeMove(move[0], cur_pos, isWhite, wasPieceCaptured, capturedPieceType)
						if kingmove == move {
							kingmove = [2]uint64{0, 0}
							castleMoveInfo["o-o"].CanCastle = true
							castleMoveInfo["o-o-o"].CanCastle = true
						}
					}
					//case missing : once king has moved, it should not be able to castle
					//implemented from user side but resetting down the children's tree not necessaerily fruitful

					if moveBeta < beta {
						beta = moveBeta
						bestMove = [2]uint64{move[0], cur_pos}
					}
					if beta <= alpha {
						to_break = true
						break
					}
				}
				if to_break {
					break
				}

			}

		}
		if bestMove == [2]uint64{0, 0} {
			bestMove = lastMove
		}
		return beta, bestMove
	}

}

func checkMoveIsCastle(move [2]uint64) (bool, string) {
	for key, value := range castleMoveInfo {
		if value.RandomMove == move {
			return true, key
		}
	}
	return false, ""
}
