package chess_test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/ekzyis/sn-chess/chess"
	"github.com/stretchr/testify/assert"
)

func init() {
	// change working directory to the root of the project
	// so assets/ can be found
	wd, _ := os.Getwd()
	os.Chdir(path.Dir(wd))
}

func TestBoardInitial(t *testing.T) {
	b := chess.NewBoard()

	assertPiece(t, b, "a1", chess.Rook, chess.Light)
	assertPiece(t, b, "b1", chess.Knight, chess.Light)
	assertPiece(t, b, "c1", chess.Bishop, chess.Light)
	assertPiece(t, b, "d1", chess.Queen, chess.Light)
	assertPiece(t, b, "e1", chess.King, chess.Light)
	assertPiece(t, b, "f1", chess.Bishop, chess.Light)
	assertPiece(t, b, "g1", chess.Knight, chess.Light)
	assertPiece(t, b, "h1", chess.Rook, chess.Light)

	assertPiece(t, b, "a2", chess.Pawn, chess.Light)
	assertPiece(t, b, "b2", chess.Pawn, chess.Light)
	assertPiece(t, b, "c2", chess.Pawn, chess.Light)
	assertPiece(t, b, "d2", chess.Pawn, chess.Light)
	assertPiece(t, b, "e2", chess.Pawn, chess.Light)
	assertPiece(t, b, "f2", chess.Pawn, chess.Light)
	assertPiece(t, b, "g2", chess.Pawn, chess.Light)
	assertPiece(t, b, "h2", chess.Pawn, chess.Light)

	assertPiece(t, b, "a8", chess.Rook, chess.Dark)
	assertPiece(t, b, "b8", chess.Knight, chess.Dark)
	assertPiece(t, b, "c8", chess.Bishop, chess.Dark)
	assertPiece(t, b, "d8", chess.Queen, chess.Dark)
	assertPiece(t, b, "e8", chess.King, chess.Dark)
	assertPiece(t, b, "f8", chess.Bishop, chess.Dark)
	assertPiece(t, b, "g8", chess.Knight, chess.Dark)
	assertPiece(t, b, "h8", chess.Rook, chess.Dark)

	assertPiece(t, b, "a7", chess.Pawn, chess.Dark)
	assertPiece(t, b, "b7", chess.Pawn, chess.Dark)
	assertPiece(t, b, "c7", chess.Pawn, chess.Dark)
	assertPiece(t, b, "d7", chess.Pawn, chess.Dark)
	assertPiece(t, b, "e7", chess.Pawn, chess.Dark)
	assertPiece(t, b, "f7", chess.Pawn, chess.Dark)
	assertPiece(t, b, "g7", chess.Pawn, chess.Dark)
	assertPiece(t, b, "h7", chess.Pawn, chess.Dark)
}

func TestBoardMovePawn(t *testing.T) {
	b := chess.NewBoard()

	b.Move("e4")

	assertNoPiece(t, b, "e2")
	assertPiece(t, b, "e4", chess.Pawn, chess.Light)

	b.Move("e5")

	assertNoPiece(t, b, "e7")
	assertPiece(t, b, "e5", chess.Pawn, chess.Dark)
}

func assertPiece(t *testing.T, b *chess.Board, position string, name chess.PieceName, color chess.Color) {
	p := b.At(position)

	c := "white"
	if color == chess.Dark {
		c = "black"
	}

	msg := fmt.Sprintf("expected %s %s at %s", c, name, position)

	assert.NotNil(t, p, msg)
	assert.Equal(t, name, p.Name, msg)
	assert.Equal(t, color, p.Color, msg)
}

func assertNoPiece(t *testing.T, b *chess.Board, position string) {
	p := b.At(position)

	assert.Nil(t, p, "expected no piece at %s", position)
}
