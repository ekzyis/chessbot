package chess_test

import (
	"testing"

	"github.com/ekzyis/chessbot/chess"
)

func TestGame001(t *testing.T) {
	t.Parallel()

	b := chess.NewBoard()

	// this used to not parse because of the + at the end
	assertParse(t, b, "d4 d5 Bf4 Nf6 e3 Ne4 Nc3 Nf2 Kxf2 e6 Qg4 Be7 Re1 O-O Kg3 Bh4+")
}
