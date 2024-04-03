package engine

// PieceType represents the type of a chess piece.
type PieceType int

// Board represents the state of the chess board.
type Board struct {
	whitePawns     uint64
	whiteKnights   uint64
	whiteBishops   uint64
	whiteRooks     uint64
	whiteQueens    uint64
	whiteKing      uint64
	whitePieces    uint64
	canWhiteCastle bool

	blackPawns     uint64
	blackKnights   uint64
	blackBishops   uint64
	blackRooks     uint64
	blackQueens    uint64
	blackKing      uint64
	blackPieces    uint64
	canBlackCastle bool

	allPieces uint64
}

const (
	bottomEdge       uint64 = 0x00000000000000FF
	bottomButOneEdge uint64 = 0x000000000000FF00

	rightEdge       uint64 = 0x0101010101010101
	rightButOneEdge uint64 = 0x0202020202020202

	leftEdge       uint64 = 0x8080808080808080
	leftButOneEdge uint64 = 0x4040404040404040

	topEdge       uint64 = 0xFF00000000000000
	topButOneEdge uint64 = 0x00FF000000000000

	diagBackRightDir    uint8 = 1 << 0
	backDir             uint8 = 1 << 1
	diagBackLeftDir     uint8 = 1 << 2
	rightDir            uint8 = 1 << 3
	leftDir             uint8 = 1 << 4
	diagForwardRightDir uint8 = 1 << 5
	forwardDir          uint8 = 1 << 6
	diagForwardLeftDir  uint8 = 1 << 7

	diagDirs     uint8 = diagBackRightDir | diagBackLeftDir | diagForwardRightDir | diagForwardLeftDir
	straightDirs uint8 = backDir | rightDir | leftDir | forwardDir
	allDirs      uint8 = diagDirs | straightDirs

	Pawn PieceType = iota + 1
	Knight
	Bishop
	Rook
	Queen
	King

	pawn_wt   = 100
	knight_wt = 320
	bishop_wt = 330
	rook_wt   = 500
	queen_wt  = 900
	king_wt   = 20000
)
