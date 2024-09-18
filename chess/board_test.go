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

func TestBoardMovePawnInvalid(t *testing.T) {
	b := chess.NewBoard()

	assertMoveError(t, b, "a5", "no pawn found that can move to a5")
	assertMoveError(t, b, "b5", "no pawn found that can move to b5")
	assertMoveError(t, b, "c5", "no pawn found that can move to c5")
	assertMoveError(t, b, "d5", "no pawn found that can move to d5")
	assertMoveError(t, b, "e5", "no pawn found that can move to e5")
	assertMoveError(t, b, "f5", "no pawn found that can move to f5")
	assertMoveError(t, b, "g5", "no pawn found that can move to g5")
	assertMoveError(t, b, "h5", "no pawn found that can move to h5")

	b.Move("d4")

	assertMoveError(t, b, "a4", "no pawn found that can move to a4")
	assertMoveError(t, b, "b4", "no pawn found that can move to b4")
	assertMoveError(t, b, "c4", "no pawn found that can move to c4")
	assertMoveError(t, b, "d4", "no pawn found that can move to d4")
	assertMoveError(t, b, "e4", "no pawn found that can move to e4")
	assertMoveError(t, b, "f4", "no pawn found that can move to f4")
	assertMoveError(t, b, "g4", "no pawn found that can move to g4")
	assertMoveError(t, b, "h4", "no pawn found that can move to h4")
}

func TestBoardMoveKnight(t *testing.T) {
	b := chess.NewBoard()

	b.Move("Nf3")

	assertPiece(t, b, "f3", chess.Knight, chess.Light)
	assertNoPiece(t, b, "g1")

	b.Move("Nh6")

	assertPiece(t, b, "h6", chess.Knight, chess.Dark)
	assertNoPiece(t, b, "g8")
}

func TestBoardMoveBishop(t *testing.T) {
	b := chess.NewBoard()

	b.Move("Bc4")

	assertMoveError(t, b, "Bc4", "no bishop found that can move to c4")

	b.Move("e3")
	b.Move("e6")

	b.Move("Bc4")

	assertPiece(t, b, "c4", chess.Bishop, chess.Light)
	assertNoPiece(t, b, "f1")

	b.Move("Bc5")

	assertPiece(t, b, "c5", chess.Bishop, chess.Dark)
	assertNoPiece(t, b, "f8")
}

func TestBoardMoveRook(t *testing.T) {
	b := chess.NewBoard()

	b.Move("Ra3")

	assertMoveError(t, b, "Ra3", "no rook found that can move to a3")

	b.Move("a4")
	b.Move("a5")

	b.Move("Ra3")

	assertPiece(t, b, "a3", chess.Rook, chess.Light)
	assertNoPiece(t, b, "a1")

	b.Move("Ra6")

	assertPiece(t, b, "a6", chess.Rook, chess.Dark)
	assertNoPiece(t, b, "a8")
}

func assertPiece(t *testing.T, b *chess.Board, position string, name chess.PieceName, color chess.Color) {
	p := b.At(position)

	c := "white"
	if color == chess.Dark {
		c = "black"
	}

	msg := fmt.Sprintf("expected %s %s at %s", c, name, position)

	if assert.NotNil(t, p, msg) {
		assert.Equal(t, name, p.Name, msg)
		assert.Equal(t, color, p.Color, msg)
	}
}

func assertNoPiece(t *testing.T, b *chess.Board, position string) {
	p := b.At(position)

	assert.Nil(t, p, "expected no piece at %s", position)
}

func assertMoveError(t *testing.T, b *chess.Board, position string, message string) {
	assert.ErrorContains(t, b.Move(position), message)
}
