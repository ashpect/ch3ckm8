package engine

import (
	"fmt"
)

var tell func(text ...string)

func init() {
	tell = tell_test // for testing
	// tell = mainTell // when not testing
}

func tell_test(text ...string) {
	// fmt.Println("Tell test:")
	fmt.Println(text)
}

func mainTell(text ...string) {

	// Just printing the text for now, send to gui later.
	toGUI := ""
	for _, t := range text {
		toGUI += t
	}

	fmt.Println("Main tell:")
	fmt.Println(toGUI)
}
