package engine

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func Uci(frGUI chan string) {

	tell("Hello from uci")

	frEng, toEng := engine()

	quit := false
	cmd := ""
	bestmove := ""

	for quit == false {
		select {
		case cmd = <-frGUI:
		case bestmove = <-frEng:
			tell(bestmove)
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
			toEng <- "stop"
		case "quit":
			quit = true
			continue
		}
	}
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
