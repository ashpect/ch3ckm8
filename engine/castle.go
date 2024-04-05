package engine

import (
	"fmt"
)

var castleMoveInfo = map[string]*CastleMoveInfo{
	"O-O":   {Positions: []string{"f1", "g1"}, CanCastle: true, RandomMove: [2]uint64{0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF}, KingMove: "Ke1-g1", RookMove: "Rh1-f1"},   //white king
	"O-O-O": {Positions: []string{"d1", "c1"}, CanCastle: true, RandomMove: [2]uint64{0xFFFFFFFFFFFFFF00, 0xFFFFFFFFFFFFFF00}, KingMove: "Ke1-c1", RookMove: "Re1-d1"}, //white queen
	"o-o":   {Positions: []string{"f8", "g8"}, CanCastle: true, RandomMove: [2]uint64{0xFFFFFFFFFFFF000, 0xFFFFFFFFFFFF0000}, KingMove: "ke8-g8", RookMove: "rh8-f8"},   //black king
	"o-o-o": {Positions: []string{"d8", "c8"}, CanCastle: true, RandomMove: [2]uint64{0xFFFFFFFFFF000000, 0xFFFFFFFFFF000000}, KingMove: "ke8-g8", RookMove: "ra8-d8"}, //black queen
}

func (b *Board) isCastleMove(move string) bool {
	_, ok := castleMoveInfo[move]
	fmt.Println(ok)
	return ok
}

// func for getting all valid castle moves available for a situation for search to assess
func (b *Board) getAllCastlingMoves(isWhite bool) [][2]uint64 {

	var castleMoves [][2]uint64

	if b.canCastleAll(isWhite) {
		return nil
	}

	for key := range castleMoveInfo {
		if b.validateCastleMove(isWhite, key) {
			castleMoves = append(castleMoves, castleMoveInfo[key].RandomMove)
		}
	}

	return castleMoves
}

// just for validating if doing the castle move is legit
// 1. the castling move is not done, and pieces are in their position (checking the struct for status)
// 2. the squares between the king and rook are empty
// 3. the squares the king moves through are not under attack
// not checking for check as it is done prior overall, can use the canCastle func and not for a move in this helper func.

func (b *Board) validateCastleMove(isWhite bool, castleMove string) bool {

	for key, value := range castleMoveInfo {

		if (key == castleMove) && value.CanCastle {
			for _, pos := range value.Positions {
				if b.validSquare(b.notationToPos(pos), isWhite) {
					return true
				}
			}
		}

	}

	return false
}

// for checking if the king is in check and if already moved
func (b *Board) canCastleAll(isWhite bool) bool {

	i := 0
	for _, value := range castleMoveInfo {
		if i <= 1 {
			if isWhite && !value.CanCastle {
				return true
			}
		}
		if i > 1 {
			if !isWhite && !value.CanCastle {
				return true
			}
		}
	}

	if b.isCheck(isWhite) {
		return false
	}
	return false
}

func (b *Board) canCastle(isWhite bool, move string) bool {

	if !castleMoveInfo[move].CanCastle {
		fmt.Println("usbfdus")
		return false
	}
	if b.isCheck(isWhite) {
		fmt.Println("yahan")
		return false
	}
	if !b.validateCastleMove(isWhite, move) {
		fmt.Println("wahan")
		return false
	}

	return true

}

func (b *Board) validSquare(pos uint64, isWhite bool) bool {
	if b.getPieceType(pos) != NoPiece {
		return false
	}
	legalMoves := b.getAllMoves(!isWhite)
	for _, move := range legalMoves {
		if move[1]&pos != 0 {
			return false
		}
	}
	return true
}

// write such that the struct gets updated whenever there is a rook move
