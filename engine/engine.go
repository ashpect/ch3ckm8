package engine

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
)

const (
	depth int = 5
)

var mainBoard Board

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
				mainBoard.Print(true)
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
				mainBoard.Print(randomBool)
				frEng <- "new board initialized, you are playing "
			case "eval":
				mainBoard.showEvalScore()
			default:
				if strings.HasPrefix(cmd, "fen ") {
					otherString := strings.TrimPrefix(cmd, "fen ")
					mainBoard = parse(otherString)
					mainBoard.Print(true)

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

	piece, initPos64, finalPos64 := b.notationToMove(move)
	colour := b.getColour(initPos64)
	_, _ = b.makeMove(initPos64, finalPos64, colour, piece)
	b.Print(colour)
	_, bestMove := b.alphaBetaMiniMax(!colour, math.Inf(-1), math.Inf(1), depth)
	var pieceType PieceType = b.getPieceType(bestMove[0])
	wasPieceCaptured, _ := b.makeMove(bestMove[0], bestMove[1], !colour, pieceType)
	b.Print(colour)
	if b.isCheckmate(!colour) {
		fmt.Println("Bot Wins!")
	} else if b.isCheckmate(colour) {
		fmt.Println("Bot Loses!")
	}

	responseMove := b.moveToNotation(colour, bestMove[0], bestMove[1], pieceType, wasPieceCaptured)
	return responseMove
}

func (b *Board) startWhite() string {
	b.Print(false)

	_, bestMove := b.alphaBetaMiniMax(true, math.Inf(-1), math.Inf(1), depth)
	responseMove := b.moveToNotation(true, bestMove[0], bestMove[1], b.getPieceType(bestMove[0]), false)
	_, _ = b.makeMove(bestMove[0], bestMove[1], true, b.getPieceType(bestMove[0]))
	b.Print(false)
	if b.isCheckmate(true) {
		fmt.Println("Bot Wins!")
	} else if b.isCheckmate(false) {
		fmt.Println("Bot Loses!")
	}

	return responseMove
}

func (b *Board) showEvalScore() {
	a, c := b.eval()
	fmt.Printf("Material score: %v", a)
	fmt.Printf("Position score: %v", c)
}
