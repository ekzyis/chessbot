package main

import (
	"github.com/ekzyis/sn-chess/chess"
)

func main() {
	var (
		b = chess.NewBoard()
	)

	b.Save("board.png")
}
