package engine

import (
	"fmt"
	"math/rand"
	"strings"
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
				mainBoard.Print(false)
				frEng <- "new board initialized, you are playing black"
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
					fmt.Println(otherString)
					mainBoard = parse(otherString)
					mainBoard.Print(true)
				}
			}
		}
	}()
	return frEng, toEng
}

func (b *Board) showEvalScore() {
	a, c := b.eval()
	fmt.Printf("Material score: %v", a)
	fmt.Printf("Position score: %v", c)
}
