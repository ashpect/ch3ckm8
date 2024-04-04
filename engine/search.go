package engine

func (b *Board) alphaBetaMiniMax(isWhite bool, depth int, alpha, beta float64) (float64, [2]uint64) {
	if depth == 0 || b.isCheckmate() {
		a, c := b.eval()
		return a + c, [2]uint64{0, 0}
	}
	if isWhite {
		maxEval := -100000.0
		bestmove := [2]uint64{0, 0}
		for _, move := range b.getAllLegalMoves(isWhite) {
			wasPieceCaptured, capturedPieceType := b.makeMove(move[0], move[1], isWhite, b.getPieceType(move[0]))
			eval, _ := b.alphaBetaMiniMax(!isWhite, depth-1, alpha, beta)
			b.unmakeMove(move[0], move[1], isWhite, wasPieceCaptured, capturedPieceType)
			if eval > maxEval {
				maxEval = eval
				bestmove = move
			}

		}
		return maxEval, bestmove
	} else {
		minEval := 100000.0
		bestmove := [2]uint64{0, 0}
		for _, move := range b.getAllLegalMoves(isWhite) {
			wasPieceCaptured, capturedPieceType := b.makeMove(move[0], move[1], isWhite, b.getPieceType(move[0]))
			eval, _ := b.alphaBetaMiniMax(!isWhite, depth-1, alpha, beta)
			b.unmakeMove(move[0], move[1], isWhite, wasPieceCaptured, capturedPieceType)
			if eval < minEval {
				beta = eval
				bestmove = move
			}
		}
		return minEval, bestmove
	}

}
