package engine

import "math"

func alphaBetaMiniMax(b Board, isWhite bool, depth int, alpha, beta float64) float64 {
	if depth == 0 {
		return b.eval()
	}
	if isWhite {
		var moves = b.getAllLegalMoves(isWhite)
		for _, move := range moves {
			for i := 0; i < 64; i++ {
				if move[1]&(1<<uint64(i)) != 0 {
					wasPieceCaptured, capturePieceType := b.makeMove(move[0], 1<<uint64(i), isWhite, b.getPieceType(move[0]))
					alpha = math.Max(alpha, alphaBetaMiniMax(b, !isWhite, depth-1, alpha, beta))
					b.unmakeMove(move[0], 1<<uint64(i), isWhite, wasPieceCaptured, capturePieceType)
					if beta <= alpha {
						break
					}
				}
			}
		}
		return alpha
	} else {
		var moves = b.getAllLegalMoves(isWhite)
		for _, move := range moves {
			for i := 0; i < 64; i++ {
				if move[1]&(1<<uint64(i)) != 0 {
					wasPieceCaptured, capturePieceType := b.makeMove(move[0], 1<<uint64(i), isWhite, b.getPieceType(move[0]))
					beta = math.Min(beta, alphaBetaMiniMax(b, !isWhite, depth-1, alpha, beta))
					b.unmakeMove(move[0], 1<<uint64(i), isWhite, wasPieceCaptured, capturePieceType)
					if beta <= alpha {
						break
					}
				}
			}
		}
		return beta
	}
}
