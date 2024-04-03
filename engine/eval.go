package engine

var (
	pawn_pos_weight = [64]float64{
		0, 0, 0, 0, 0, 0, 0, 0,
		50, 50, 50, 50, 50, 50, 50, 50,
		10, 10, 20, 30, 30, 20, 10, 10,
		5, 5, 10, 25, 25, 10, 5, 5,
		0, 0, 0, 20, 20, 0, 0, 0,
		5, -5, -10, 0, 0, -10, -5, 5,
		5, 10, 10, -20, -20, 10, 10, 5,
		0, 0, 0, 0, 0, 0, 0, 0}

	pawn_pos_weight_eg = [64]float64{
		0, 0, 0, 0, 0, 0, 0, 0,
		80, 80, 80, 80, 80, 80, 80, 80,
		50, 50, 50, 50, 50, 50, 50, 50,
		30, 30, 30, 30, 30, 30, 30, 30,
		20, 20, 20, 20, 20, 20, 20, 20,
		10, 10, 10, 10, 10, 10, 10, 10,
		10, 10, 10, 10, 10, 10, 10, 10,
		0, 0, 0, 0, 0, 0, 0, 0}

	knight_pos_weight = [64]float64{-50, -40, -30, -30, -30, -30, -40, -50,
		-40, -20, 0, 0, 0, 0, -20, -40,
		-30, 0, 10, 15, 15, 10, 0, -30,
		-30, 5, 15, 20, 20, 15, 5, -30,
		-30, 0, 15, 20, 20, 15, 0, -30,
		-30, 5, 10, 15, 15, 10, 5, -30,
		-40, -20, 0, 5, 5, 0, -20, -40,
		-50, -40, -30, -30, -30, -30, -40, -50}

	bishop_pos_weight = [64]float64{-20, -10, -10, -10, -10, -10, -10, -20,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-10, 0, 5, 10, 10, 5, 0, -10,
		-10, 5, 5, 10, 10, 5, 5, -10,
		-10, 0, 10, 10, 10, 10, 0, -10,
		-10, 10, 10, 10, 10, 10, 10, -10,
		-10, 5, 0, 0, 0, 0, 5, -10,
		-20, -10, -10, -10, -10, -10, -10, -20}

	rook_pos_weight = [64]float64{0, 0, 0, 0, 0, 0, 0, 0,
		5, 10, 10, 10, 10, 10, 10, 5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		0, 0, 0, 5, 5, 0, 0, 0}

	queen_pos_weight = [64]float64{-20, -10, -10, -5, -5, -10, -10, -20,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-10, 0, 5, 5, 5, 5, 0, -10,
		-5, 0, 5, 5, 5, 5, 0, -5,
		0, 0, 5, 5, 5, 5, 0, -5,
		-10, 5, 5, 5, 5, 5, 0, -10,
		-10, 0, 5, 0, 0, 0, 0, -10,
		-20, -10, -10, -5, -5, -10, -10, -20}

	king_pos_weight = [64]float64{-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-20, -30, -30, -40, -40, -30, -30, -20,
		-10, -20, -20, -20, -20, -20, -20, -10,
		20, 20, 0, 0, 0, 0, 20, 20,
		20, 30, 10, 0, 0, 10, 30, 20}

	king_pos_weight_eg = [64]float64{-50, -40, -30, -20, -20, -30, -40, -50,
		-30, -20, -10, 0, 0, -10, -20, -30,
		-30, -10, 20, 30, 30, 20, -10, -30,
		-30, -10, 30, 40, 40, 30, -10, -30,
		-30, -10, 30, 40, 40, 30, -10, -30,
		-30, -10, 20, 30, 30, 20, -10, -30,
		-30, -30, 0, 0, 0, 0, -30, -30,
		-50, -30, -30, -30, -30, -30, -30, -50}
)

func (b *Board) eval() float64 {
	var endgameT float64 = 0
	var materialScore float64 = b.evalMaterialValues()
	var pstScore float64 = b.evalPieceSquareTables(endgameT)
	return materialScore + pstScore
}

func (b *Board) evalMaterialValues() float64 {
	var score float64 = 0
	for i := 0; i < 64; i++ {
		var cur_pos uint64 = 1 << uint64(i)
		if b.whitePawns&cur_pos != 0 {
			score += pawn_wt
		} else if b.whiteKnights&cur_pos != 0 {
			score += knight_wt
		} else if b.whiteBishops&cur_pos != 0 {
			score += bishop_wt
		} else if b.whiteRooks&cur_pos != 0 {
			score += rook_wt
		} else if b.whiteQueens&cur_pos != 0 {
			score += queen_wt
		} else if b.whiteKing&cur_pos != 0 {
			score += king_wt
		} else if b.blackPawns&cur_pos != 0 {
			score -= pawn_wt
		} else if b.blackKnights&cur_pos != 0 {
			score -= knight_wt
		} else if b.blackBishops&cur_pos != 0 {
			score -= bishop_wt
		} else if b.blackRooks&cur_pos != 0 {
			score -= rook_wt
		} else if b.blackQueens&cur_pos != 0 {
			score -= queen_wt
		} else if b.blackKing&cur_pos != 0 {
			score -= king_wt
		}
	}
	return score
}
func (b *Board) evalPieceSquareTables(endgameT float64) float64 {
	var score float64 = 0
	for i := 0; i < 64; i++ {
		var cur_pos uint64 = 1 << uint64(i)
		if b.whitePawns&cur_pos != 0 {
			score += pawn_pos_weight[63-i]*(1-endgameT) + pawn_pos_weight_eg[63-i]*endgameT
		} else if b.whiteKnights&cur_pos != 0 {
			score += knight_pos_weight[63-i]
		} else if b.whiteBishops&cur_pos != 0 {
			score += bishop_pos_weight[63-i]
		} else if b.whiteRooks&cur_pos != 0 {
			score += rook_pos_weight[63-i]
		} else if b.whiteQueens&cur_pos != 0 {
			score += queen_pos_weight[63-i]
		} else if b.whiteKing&cur_pos != 0 {
			score += king_pos_weight[63-i]*(1-endgameT) + king_pos_weight_eg[63-i]*endgameT
		} else if b.blackPawns&cur_pos != 0 {
			score -= pawn_pos_weight[i]*(1-endgameT) + pawn_pos_weight_eg[i]*endgameT
		} else if b.blackKnights&cur_pos != 0 {
			score -= knight_pos_weight[i]
		} else if b.blackBishops&cur_pos != 0 {
			score -= bishop_pos_weight[i]
		} else if b.blackRooks&cur_pos != 0 {
			score -= rook_pos_weight[i]
		} else if b.blackQueens&cur_pos != 0 {
			score -= queen_pos_weight[i]
		} else if b.blackKing&cur_pos != 0 {
			score -= king_pos_weight[i]*(1-endgameT) + king_pos_weight_eg[i]*endgameT
		}
	}
	return score
}
