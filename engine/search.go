package engine

import (
	"fmt"
	"math"
)

func (b *Board) alphaBetaMiniMax(isWhite bool, depth int) (float64, [2]uint64) {
	if depth == 0 {
		a, c := b.eval()
		return a + c, [2]uint64{0, 0}
	}
	if isWhite {
		maxEval := math.Inf(-1)
		bestMove := [2]uint64{0, 0}
		for _, move := range b.getAllLegalMoves(isWhite) {
			for i := 0; i < 64; i++ {
				var cur_pos uint64 = 1 << uint64(i)
				if move[1]&cur_pos != 0 {
					wasPieceCaptured, capturedPieceType := b.makeMove(move[0], cur_pos, isWhite, b.getPieceType(move[0]))
					eval, _ := b.alphaBetaMiniMax(!isWhite, depth-1)
					b.unmakeMove(move[0], cur_pos, isWhite, wasPieceCaptured, capturedPieceType)
					if eval > maxEval {
						maxEval = eval
						bestMove = [2]uint64{move[0], cur_pos}
					}
				}
			}

		}
		return maxEval, bestMove
	} else {
		minEval := math.Inf(1)
		bestMove := [2]uint64{0, 0}
		for _, move := range b.getAllLegalMoves(isWhite) {
			for i := 0; i < 64; i++ {
				var cur_pos uint64 = 1 << uint64(i)
				if move[1]&cur_pos != 0 {
					wasPieceCaptured, capturedPieceType := b.makeMove(move[0], cur_pos, isWhite, b.getPieceType(move[0]))
					eval, _ := b.alphaBetaMiniMax(!isWhite, depth-1)
					b.unmakeMove(move[0], cur_pos, isWhite, wasPieceCaptured, capturedPieceType)
					if eval < minEval {
						minEval = eval
						bestMove = [2]uint64{move[0], cur_pos}

					}
				}
			}

		}
		return minEval, bestMove
	}

}

func (b *Board) searchToMove(isWhite bool, moves [2]uint64) string {

	initPiece := b.piecePosToNotation(moves[0])
	initPos := b.posToNotation(moves[0])
	finalPos := b.posToNotation(moves[1])
	if !isWhite {
		initPos = flipNotation(initPos)
		finalPos = flipNotation(finalPos)
	}

	return initPiece + initPos + finalPos
}

func (b *Board) moveToSearch(move string) (PieceType, uint64, uint64) {
	initPos := move[2:4]
	finalPos := move[4:6]
	initPiece := move[0:1]

	fmt.Println(initPiece, initPos, finalPos)

	piece := b.notationToPieceType(initPiece)
	initPos64 := b.notationToPos(initPos)
	finalPos64 := b.notationToPos(finalPos)

	return piece, initPos64, finalPos64
}
