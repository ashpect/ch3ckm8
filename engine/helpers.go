package engine

import "strings"

func (b *Board) getPieceNotation(pos uint64) string {

	pieceType := b.getPieceType(pos)
	isWhite := b.getColour(pos)
	var base_string string = ""
	switch pieceType {
	case Pawn:
		base_string = ""
	case Knight:
		base_string = "N"
	case Bishop:
		base_string = "B"
	case Rook:
		base_string = "R"
	case Queen:
		base_string = "Q"
	case King:
		base_string = "K"
	}
	if !isWhite {
		base_string = strings.ToLower(base_string)
	}
	return base_string
}

func (b *Board) getColour(pos uint64) bool {
	return pos&b.whitePieces != 0
}

func (b *Board) getPosNotation(pos uint64) string {
	var posNotation string = ""
	var i uint64
	for i = 0x8000000000000000; i > 0; i >>= 1 {
		if pos&i != 0 {
			posNotation = string(rune(97+(63-i)%8)) + string(rune(49+int((63-i)/8)))
			break
		}
	}
	return posNotation
}
