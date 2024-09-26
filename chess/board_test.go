package chess_test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/ekzyis/chessbot/chess"
	"github.com/stretchr/testify/assert"
)

func init() {
	// change working directory to the root of the project
	// so assets/ can be found
	wd, _ := os.Getwd()
	os.Chdir(path.Dir(wd))
}

func TestBoardInitial(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

	b := chess.NewBoard()

	assertParse(t, b, "e4")

	assertNoPiece(t, b, "e2")
	assertPiece(t, b, "e4", chess.Pawn, chess.Light)

	assertParse(t, b, "e5")

	assertNoPiece(t, b, "e7")
	assertPiece(t, b, "e5", chess.Pawn, chess.Dark)
}

func TestBoardMovePawnInvalid(t *testing.T) {
	t.Parallel()

	b := chess.NewBoard()

	assertMoveError(t, b, "a5", "no pawn found that can move to a5")
	assertMoveError(t, b, "b5", "no pawn found that can move to b5")
	assertMoveError(t, b, "c5", "no pawn found that can move to c5")
	assertMoveError(t, b, "d5", "no pawn found that can move to d5")
	assertMoveError(t, b, "e5", "no pawn found that can move to e5")
	assertMoveError(t, b, "f5", "no pawn found that can move to f5")
	assertMoveError(t, b, "g5", "no pawn found that can move to g5")
	assertMoveError(t, b, "h5", "no pawn found that can move to h5")

	assertParse(t, b, "d4")

	assertMoveError(t, b, "a4", "no pawn found that can move to a4")
	assertMoveError(t, b, "b4", "no pawn found that can move to b4")
	assertMoveError(t, b, "c4", "no pawn found that can move to c4")
	assertMoveError(t, b, "d4", "no pawn found that can move to d4")
	assertMoveError(t, b, "e4", "no pawn found that can move to e4")
	assertMoveError(t, b, "f4", "no pawn found that can move to f4")
	assertMoveError(t, b, "g4", "no pawn found that can move to g4")
	assertMoveError(t, b, "h4", "no pawn found that can move to h4")
}

func TestBoardMovePawnCapture(t *testing.T) {
	t.Parallel()

	b := chess.NewBoard()

	assertParse(t, b, "e4 d5 exd5")

	assertNoPiece(t, b, "e4")
	assertPiece(t, b, "d5", chess.Pawn, chess.Light)

	// test ambiguous capture

	b = chess.NewBoard()

	assertParse(t, b, "c4 d5 e4 e5 exd5")

	assertNoPiece(t, b, "e4")
	assertPiece(t, b, "d5", chess.Pawn, chess.Light)
	assertPiece(t, b, "c4", chess.Pawn, chess.Light)

	b = chess.NewBoard()

	assertParse(t, b, "c4 d5 e4 e5 cxd5")

	assertNoPiece(t, b, "c4")
	assertPiece(t, b, "d5", chess.Pawn, chess.Light)
	assertPiece(t, b, "e4", chess.Pawn, chess.Light)
}

func TestBoardMoveKnight(t *testing.T) {
	t.Parallel()

	b := chess.NewBoard()

	assertParse(t, b, "Nf3")
	assertPiece(t, b, "f3", chess.Knight, chess.Light)
	assertNoPiece(t, b, "g1")

	assertParse(t, b, "Nh6")
	assertPiece(t, b, "h6", chess.Knight, chess.Dark)
	assertNoPiece(t, b, "g8")

	assertParse(t, b, "Nc3")
	assertPiece(t, b, "c3", chess.Knight, chess.Light)
	assertNoPiece(t, b, "b1")

	assertParse(t, b, "Na6")
	assertPiece(t, b, "a6", chess.Knight, chess.Dark)
	assertNoPiece(t, b, "b8")

	assertParse(t, b, "Nh4")
	assertPiece(t, b, "h4", chess.Knight, chess.Light)
	assertNoPiece(t, b, "f3")

	assertParse(t, b, "Nf5")
	assertPiece(t, b, "f5", chess.Knight, chess.Dark)
	assertNoPiece(t, b, "h6")

	assertParse(t, b, "Na4")
	assertPiece(t, b, "a4", chess.Knight, chess.Light)
	assertNoPiece(t, b, "c3")

	assertParse(t, b, "Nc5")
	assertPiece(t, b, "c5", chess.Knight, chess.Dark)
	assertNoPiece(t, b, "a6")
}

func TestBoardMoveKnightInvalid(t *testing.T) {
	t.Parallel()

	b := chess.NewBoard()

	// out of reach
	assertMoveError(t, b, "Ng3", "no knight found that can move to g3")
	assertMoveError(t, b, "Nb3", "no knight found that can move to b3")

	// blocked by own piece
	assertMoveError(t, b, "Ng1", "g1 blocked by white knight")
	assertMoveError(t, b, "Nd2", "d2 blocked by white pawn")
	assertMoveError(t, b, "Ne2", "e2 blocked by white pawn")

	assertParse(t, b, "Nf3")

	assertMoveError(t, b, "Ng6", "no knight found that can move to g6")
	assertMoveError(t, b, "Nb6", "no knight found that can move to b6")

	assertMoveError(t, b, "Ne7", "e7 blocked by black pawn")
	assertMoveError(t, b, "Nd7", "d7 blocked by black pawn")
}

func TestBoardMoveKnightCapture(t *testing.T) {
	t.Parallel()

	b := chess.NewBoard()

	assertParse(t, b, "e4 Nf6 d4 Nxe4")

	assertPiece(t, b, "e4", chess.Knight, chess.Dark)
	assertPiece(t, b, "d4", chess.Pawn, chess.Light)
	assertNoPiece(t, b, "g8")
	assertNoPiece(t, b, "e2")
	assertNoPiece(t, b, "d2")

	// test ambiguous capture

	b = chess.NewBoard()

	assertParse(t, b, "e4 e5 Nf3 d6 Nc3 d5 Nb5 d4 Nbxd4")

	assertPiece(t, b, "e4", chess.Pawn, chess.Light)
	assertPiece(t, b, "e5", chess.Pawn, chess.Dark)
	assertPiece(t, b, "d4", chess.Knight, chess.Light)
	assertPiece(t, b, "f3", chess.Knight, chess.Light)
	assertNoPiece(t, b, "g1")
	assertNoPiece(t, b, "b1")
	assertNoPiece(t, b, "c3")
	assertNoPiece(t, b, "d6")
	assertNoPiece(t, b, "d5")
	assertNoPiece(t, b, "b5")

	b = chess.NewBoard()

	assertParse(t, b, "e4 e5 Nf3 d6 Nc3 d5 Nb5 d4 Nfxd4")

	assertPiece(t, b, "e4", chess.Pawn, chess.Light)
	assertPiece(t, b, "e5", chess.Pawn, chess.Dark)
	assertPiece(t, b, "d4", chess.Knight, chess.Light)
	assertPiece(t, b, "b5", chess.Knight, chess.Light)
	assertNoPiece(t, b, "g1")
	assertNoPiece(t, b, "b1")
	assertNoPiece(t, b, "c3")
	assertNoPiece(t, b, "d6")
	assertNoPiece(t, b, "d5")
	assertNoPiece(t, b, "f3")
}

func TestBoardMoveBishop(t *testing.T) {
	t.Parallel()

	b := chess.NewBoard()

	assertMoveError(t, b, "Bc4", "no bishop found that can move to c4")

	assertParse(t, b, "e3")
	assertParse(t, b, "e6")

	assertParse(t, b, "Bc4")

	assertPiece(t, b, "c4", chess.Bishop, chess.Light)
	assertNoPiece(t, b, "f1")

	assertParse(t, b, "Bc5")

	assertPiece(t, b, "c5", chess.Bishop, chess.Dark)
	assertNoPiece(t, b, "f8")
}

func TestBoardMoveBishopInvalid(t *testing.T) {
	t.Parallel()

	b := chess.NewBoard()

	assertMoveError(t, b, "Bc3", "no bishop found that can move to c3")
	assertMoveError(t, b, "Bc2", "c2 blocked by white pawn")
	assertMoveError(t, b, "Bb2", "b2 blocked by white pawn")
}

func TestBoardMoveRook(t *testing.T) {
	t.Parallel()

	b := chess.NewBoard()

	assertMoveError(t, b, "Ra3", "no rook found that can move to a3")

	assertParse(t, b, "a4")
	assertParse(t, b, "a5")

	assertParse(t, b, "Ra3")

	assertPiece(t, b, "a3", chess.Rook, chess.Light)
	assertNoPiece(t, b, "a1")

	assertParse(t, b, "Ra6")

	assertPiece(t, b, "a6", chess.Rook, chess.Dark)
	assertNoPiece(t, b, "a8")
}

func TestBoardMoveRookInvalid(t *testing.T) {
	t.Parallel()

	b := chess.NewBoard()

	assertMoveError(t, b, "Rb2", "b2 blocked by white pawn")
	assertMoveError(t, b, "Rb1", "b1 blocked by white knight")
	assertMoveError(t, b, "Ra2", "a2 blocked by white pawn")
	assertMoveError(t, b, "Ra3", "no rook found that can move to a3")
}

func TestBoardMoveQueen(t *testing.T) {
	t.Parallel()

	b := chess.NewBoard()

	assertMoveError(t, b, "Qd3", "no queen found that can move to d3")

	assertParse(t, b, "d4")
	assertParse(t, b, "d5")

	assertParse(t, b, "Qd3")

	assertPiece(t, b, "d3", chess.Queen, chess.Light)
	assertNoPiece(t, b, "d1")

	assertParse(t, b, "Qd6")

	assertPiece(t, b, "d6", chess.Queen, chess.Dark)
	assertNoPiece(t, b, "d8")
}

func TestBoardMoveQueenInvalid(t *testing.T) {
	t.Parallel()

	b := chess.NewBoard()

	assertMoveError(t, b, "Qd2", "d2 blocked by white pawn")
	assertMoveError(t, b, "Qd1", "d1 blocked by white queen")
	assertMoveError(t, b, "Qe1", "e1 blocked by white king")
	assertMoveError(t, b, "Qc1", "c1 blocked by white bishop")
	assertMoveError(t, b, "Qd3", "no queen found that can move to d3")
}

func TestBoardMoveKing(t *testing.T) {
	t.Parallel()

	b := chess.NewBoard()

	assertParse(t, b, "e4 e5 Ke2 Ke7 Kf3 Kd6 Kg3 Kc6")
	assertNoPiece(t, b, "e1")
	assertNoPiece(t, b, "e8")
	assertNoPiece(t, b, "e2")
	assertNoPiece(t, b, "e7")
	assertNoPiece(t, b, "f3")
	assertNoPiece(t, b, "d6")

	assertPiece(t, b, "g3", chess.King, chess.Light)
	assertPiece(t, b, "c6", chess.King, chess.Dark)
}

func TestBoardMoveKingInvalid(t *testing.T) {
	t.Parallel()

	b := chess.NewBoard()

	assertMoveError(t, b, "Ke1", "e1 blocked by white king")
	assertMoveError(t, b, "Ke2", "e2 blocked by white pawn")
	assertMoveError(t, b, "Ke3", "no king found that can move to e3")
}

func assertParse(t *testing.T, b *chess.Board, moves string) {
	assert.NoError(t, b.Parse(moves))
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
