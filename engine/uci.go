package engine

import (
	"bufio"
	"io"
	"os"
	"strings"
)

var saveBm = ""

func Uci(frGUI chan string) {

	tell("Hello from uci")

	frEng, toEng := engine()

	bInfinite := false
	quit := false
	cmd := ""
	bestmove := ""

	for quit == false {
		select {
		case cmd = <-frGUI:
		case bestmove = <-frEng:
			handleBm(bestmove, &bInfinite)
			continue
		}
		switch cmd {
		case "uci":
			//check if uci on and details and etc.
			handleUci()
		case "debug on":
			// to add a debug mode later
			handleDebugOn()
		case "isready":
			// just send ready to test
			handleIsReady()
		case "stop":
			handleStop(toEng, &bInfinite)
		case "test":
			handleTest(toEng)
		case "newgame w":
			toEng <- "w"
		case "newgame b":
			toEng <- "b"
		case "newgame":
			toEng <- "random"
		case "eval":
			toEng <- "eval"
		default:
			if strings.HasPrefix(cmd, "position ") {
				otherString := strings.TrimPrefix(cmd, "position ")
				handlePosition(toEng, otherString)
			} else if strings.HasPrefix(cmd, "move ") {
				otherString := strings.TrimPrefix(cmd, "move ")
				handleMove(toEng, otherString)
			}
		case "quit":
			quit = true
			continue
		}
	}
}

func handleMove(toEng chan string, otherString string) {
	toEng <- "move " + otherString
}

func handlePosition(toEng chan string, otherString string) {
	toEng <- "fen " + otherString
}

func handleTest(toEng chan string) {
	toEng <- "test"
	tell("test done")
}

func handleDebugOn() {
	tell("info string debug on")
}

func handleUci() {
	tell("id name Ashish")
	tell("id author Ashish")
	tell("uciok")
}

func handleIsReady() {
	tell("readyok")
}

func handleStop(toEng chan string, bInfinite *bool) {

	// stop not valid if engine is not in infinite mode
	if *bInfinite {
		// basically if i start engine as go infinite, then i need to send stop command to engine, hence if binfinite is true, then send stop command, which was set true in go command
		// save best move and don't sent untill stop is given if the engine is in infinite mode, but had finite choices, like mate in one and exited
		// engine should be clean
		// from the infinite loop
		if saveBm != "" {
			tell(saveBm)
			saveBm = ""
		}

		toEng <- "stop"
		*bInfinite = false

	}
}

func handleBm(bestmove string, binfinite *bool) {
	if *binfinite {
		saveBm = bestmove
	} else {
		tell(bestmove)
	}
}

// go routine waits for commands and sends to uci : standard way to use anonymous func as go routines
// even when we leave the input function, the go routine will continue to run in the background
func Input() chan string {
	line := make(chan string)
	go func() {
		var reader *bufio.Reader
		//A buffered reader that reads from the standard input (os.Stdin).
		reader = bufio.NewReader(os.Stdin)
		for {
			//To read a line of text from the standard input until a newline character ('\n') is encountered
			text, err := reader.ReadString('\n')
			text = strings.TrimSpace(text)
			if err != io.EOF && len(text) > 0 {
				// send text through channel
				line <- text
			}
		}
	}()
	return line
}
