package engine

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
				b.Print(true)
				frEng <- "new board initialized, you are playing white"
			case "b":
				var b Board
				b.Initialize()
				b.Print(false)
				frEng <- "new board initialized, you are playing black"
			case "random":
				var b Board
				b.Initialize()
				// randomNumber := rand.Intn(2)
				// randomBool := randomNumber == 1
				b.Print(false)
				frEng <- "new board initialized, you are playing "
			}
		}
	}()
	return frEng, toEng
}
