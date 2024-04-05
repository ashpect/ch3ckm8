package engine

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
)

const (
	searchDepth int = 5
)

var mainBoard Board

// better design to also store the result of getalllegalmoves and keep updating as this func is called many times.
var moveHistory []string

func engine() (frEng chan string, toEng chan string) {
	tell("Hello from engine")
	frEng = make(chan string)
	toEng = make(chan string)
	go func() {
		for cmd := range toEng {
			switch cmd {
			case "stop":
				// stop here but keep channel active as engine can stop in case of 1 move away checkmate type conditions and infinite depth initiated by go inifinite
				frEng <- "stop as text recieves"
			case "quit":
				break
			case "test":
				test()
				frEng <- "test done"
			case "w":
				var b Board
				b.Initialize()
				mainBoard = b
				mainBoard.PrintBoard(true, 0)
				frEng <- "new board initialized, you are playing white"
			case "b":
				var b Board
				b.Initialize()
				mainBoard = b
				frEng <- "new board initialized, you are playing black"
				mainBoard.startWhite()
			case "random":
				var b Board
				b.Initialize()
				randomNumber := rand.Intn(2)
				randomBool := randomNumber == 1
				mainBoard = b
				mainBoard.PrintBoard(randomBool, 0)
				frEng <- "new board initialized, you are playing "
			case "eval":
				mainBoard.showEvalScore()
			default:
				if strings.HasPrefix(cmd, "fen ") {
					otherString := strings.TrimPrefix(cmd, "fen ")
					mainBoard = parse(otherString)
					mainBoard.PrintBoard(true, 0)

				} else if strings.HasPrefix(cmd, "move ") {
					otherString := strings.TrimPrefix(cmd, "move ")
					responseMove := mainBoard.handleMove(otherString)
					frEng <- responseMove
				}
			}

		}
	}()
	return frEng, toEng
}

func (b *Board) handleMove(move string) string {

	if b.isCastleMove(move) {
		submove := castleMoveInfo[move].KingMove[1:3]
		colour := b.getColour(b.notationToPos(submove))
		fmt.Println(colour, submove)
		fmt.Println(castleMoveInfo[move].CanCastle)
		if b.canCastle(colour, move) {
			responseMove := b.handleCastleMove(move)
			castleMoveInfo[move].CanCastle = false
			return responseMove
		} else {
			return "Illegal Move"
		}
	}

	piece, initPos64, finalPos64 := b.notationToMove(move)
	colour := b.getColour(initPos64)

	if !b.isMoveLegal(colour, initPos64, finalPos64) {
		return "Illegal Move"
	}
	_, _ = b.makeMove(initPos64, finalPos64, colour, piece)
	updateCastleMoveInfo(move)
	b.PrintBoard(colour, finalPos64)

	responseMove := b.getResponseMove(colour)
	updateCastleMoveInfo(responseMove)
	// Edge case : when engine plays king in one of its position children and goes back, he can castle cause we are updating here. So, handle that case.
	// First of all enable for engine the castle move.
	return responseMove
}

func (b *Board) getResponseMove(colour bool) string {
	_, bestMove := b.alphaBetaMiniMax(!colour, math.Inf(-1), math.Inf(1), searchDepth)
	var pieceType PieceType = b.getPieceType(bestMove[0])
	wasPieceCaptured, _ := b.makeMove(bestMove[0], bestMove[1], !colour, pieceType)

	b.PrintBoard(colour, bestMove[1])

	if b.isCheckmate(colour) {
		fmt.Println("Bot Wins!")
	} else if b.isCheckmate(!colour) {
		fmt.Println("Bot Loses!")
	}

	responseMove := b.moveToNotation(colour, bestMove[0], bestMove[1], pieceType, wasPieceCaptured)
	return responseMove
}

func (b *Board) handleCastleMove(move string) string {

	kingpiece, kinginitPos64, kingfinalPos64 := b.notationToMove(castleMoveInfo[move].KingMove)
	colour := b.getColour(kinginitPos64)

	_, _ = b.makeMove(kinginitPos64, kingfinalPos64, colour, kingpiece)
	moveHistory = append(moveHistory, move+" ")

	rookpiece, rookinitPos64, rookfinalPos64 := b.notationToMove(castleMoveInfo[move].RookMove)

	_, _ = b.makeMove(rookinitPos64, rookfinalPos64, colour, rookpiece)
	moveHistory = append(moveHistory, move+" ")

	b.PrintBoard(colour, rookfinalPos64)

	responseMove := b.getResponseMove(colour)
	return responseMove
}

func updateCastleMoveInfo(move string) {

	if move[0] == 'K' {
		castleMoveInfo["O-O"].CanCastle = false
		castleMoveInfo["O-O-O"].CanCastle = false
	}
	if move[0] == 'k' {
		castleMoveInfo["o-o"].CanCastle = false
		castleMoveInfo["o-o-o"].CanCastle = false
	}

}

// func (b *Board) makeUserMove(move string) bool {
// 	piece, initPos64, finalPos64 := b.moveToSearch(move)
// 	colour := b.getColour(initPos64)
// 	_, _ = b.makeMove(initPos64, finalPos64, colour, piece)
// 	b.Print(colour)
// 	return colour
// }

// func (b *Board) makeBestMove(isWhite bool) string {
// 	_, bestMove := b.alphaBetaMiniMax(!isWhite, math.Inf(-1), math.Inf(1), depth)
// 	responseMove := b.moveToReadableMove(isWhite, bestMove)
// 	_, _ = b.makeMove(bestMove[0], bestMove[1], !isWhite, b.getPieceType(bestMove[0]))
// 	return responseMove
// }

// func (b *Board) startWhite() {
// 	b.Print(false)
// 	b.makeBestMove(false)
// 	b.Print(false)
// }

func (b *Board) startWhite() string {
	b.PrintBoard(false, 0)
	responseMove := b.getResponseMove(false)
	return responseMove
}

func (b *Board) showEvalScore() {
	materialScore, positonalScore := b.eval()
	fmt.Printf("Material score: %v\n", materialScore)
	fmt.Printf("Position score: %v\n", positonalScore)
	fmt.Printf("Mobility score: %v\n", b.evalMobility())
}
