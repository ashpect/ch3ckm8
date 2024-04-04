package engine

import (
	"fmt"
)

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
					wasPieceCaptured, capturedPieceType := b.makeMove(move[0], cur_pos, isWhite, b.getPieceType(move[0]))
					moveAlpha, _ := b.alphaBetaMiniMax(!isWhite, alpha, beta, depth-1)
					b.unmakeMove(move[0], cur_pos, isWhite, wasPieceCaptured, capturedPieceType)
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
			for i := 0; i < 64; i++ {
				var to_break bool = false
				var cur_pos uint64 = 1 << uint64(i)
				if move[1]&cur_pos != 0 {
					lastMove = [2]uint64{move[0], cur_pos}
					wasPieceCaptured, capturedPieceType := b.makeMove(move[0], cur_pos, isWhite, b.getPieceType(move[0]))
					moveBeta, _ := b.alphaBetaMiniMax(!isWhite, alpha, beta, depth-1)
					b.unmakeMove(move[0], cur_pos, isWhite, wasPieceCaptured, capturedPieceType)
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
