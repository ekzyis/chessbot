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

	// this used to crash because of an out of bounds error (position was parsed as '6')
	assertMoveError(t, b, "Qex6", "square does not exist")
}
