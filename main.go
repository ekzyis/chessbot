package main

import (
	"github.com/ekzyis/chessbot/chess"
)

func main() {
	var (
		b = chess.NewBoard()
	)

	b.Save("board.png")
}
