package engine

import (
	"math"
)

func (b *Board) alphaBetaMiniMax(isWhite bool, depth int) (float64, [2]uint64) {
	if depth == 0 || b.isCheckmate() {
		a, c := b.eval()
		return a + c, [2]uint64{0, 0}
	}
	if isWhite {
		maxEval := math.Inf(-1)
		bestMove := [2]uint64{0, 0}
		for _, move := range b.getAllLegalMoves(isWhite) {
			wasPieceCaptured, capturedPieceType := b.makeMove(move[0], move[1], isWhite, b.getPieceType(move[0]))
			eval, _ := b.alphaBetaMiniMax(!isWhite, depth-1)
			b.unmakeMove(move[0], move[1], isWhite, wasPieceCaptured, capturedPieceType)
			if eval > maxEval {
				maxEval = eval
				bestMove = move
			}

		}
		return maxEval, bestMove
	} else {
		minEval := math.Inf(1)
		bestMove := [2]uint64{0, 0}
		for _, move := range b.getAllLegalMoves(isWhite) {
			wasPieceCaptured, capturedPieceType := b.makeMove(move[0], move[1], isWhite, b.getPieceType(move[0]))
			eval, _ := b.alphaBetaMiniMax(!isWhite, depth-1)
			b.unmakeMove(move[0], move[1], isWhite, wasPieceCaptured, capturedPieceType)
			if eval < minEval {
				minEval = eval
				bestMove = move
			}
		}
		return minEval, bestMove
	}

}

func (b *Board) searchToMove(moves [2]uint64) string {
	initPiece := b.piecePosToNotation(moves[0])
	initPos := posToNotation(moves[0])
	finalPos := posToNotation(moves[1])

	return initPiece + initPos + finalPos
}
