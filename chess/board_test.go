package chess_test

import (
	"fmt"
	"os"
	"path"
	"strings"
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

func TestBoardPawnPromotion(t *testing.T) {
	t.Parallel()

	b := chess.NewBoard()
	b.Parse("a4 e6 a5 e5 a6 e4 axb7 e3 bxa8=Q")

	assertPiece(t, b, "a8", chess.Queen, chess.Light)
	assertNoPiece(t, b, "b7")

	b = chess.NewBoard()
	b.Parse("a4 e6 a5 e5 a6 e4 axb7 e3 bxa8=R")

	assertPiece(t, b, "a8", chess.Rook, chess.Light)
	assertNoPiece(t, b, "b7")

	b = chess.NewBoard()
	b.Parse("a4 e6 a5 e5 a6 e4 axb7 e3 bxa8=B")

	assertPiece(t, b, "a8", chess.Bishop, chess.Light)
	assertNoPiece(t, b, "b7")

	b = chess.NewBoard()
	b.Parse("a4 e6 a5 e5 a6 e4 axb7 e3 bxa8=N")

	assertPiece(t, b, "a8", chess.Knight, chess.Light)
	assertNoPiece(t, b, "b7")

	b = chess.NewBoard()
	b.Parse("a4 e6 a5 e5 a6 e4 axb7 e3")

	assertMoveError(t, b, "bxa8=K", "invalid promotion: K")
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

	// ambiguous moves

	b = chess.NewBoard()
	assertParse(t, b, "e4 e5 Nf3 d6 Nc3 d5 Nb5 d4")
	assertMoveError(t, b, "Nxd4", "move ambiguous: 2 knights can move to d4")
	assertMoveError(t, b, "N4xd4", "no knight found that can move to d4")
	// disambiguate via file
	assertParse(t, b, "Nfxd4")

	b = chess.NewBoard()
	assertParse(t, b, "e4 e5 Nf3 d6 Nc3 d5 Nb5 d4")
	// disambiguate via rank
	assertParse(t, b, "N3xd4")
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

	// path blocked by white pawn at e2
	assertMoveError(t, b, "Bd3", "no bishop found that can move to d3")
}

func TestBoardMoveBishopCapture(t *testing.T) {
	t.Parallel()

	b := chess.NewBoard()

	assertParse(t, b, "e4 e5 Bc4 d5 Bxd5")

	assertPiece(t, b, "e4", chess.Pawn, chess.Light)
	assertPiece(t, b, "e5", chess.Pawn, chess.Dark)
	assertPiece(t, b, "d5", chess.Bishop, chess.Light)
	assertNoPiece(t, b, "f1")
	assertNoPiece(t, b, "e2")
	assertNoPiece(t, b, "e7")
	assertNoPiece(t, b, "d7")

	// bishop captures are never ambiguous because they can only move on same color
}

func TestBoardMoveRook(t *testing.T) {
	t.Parallel()

	b := chess.NewBoard()

	assertMoveError(t, b, "Ra3", "no rook found that can move to a3")

	assertParse(t, b, "a4 a5 Ra3")

	assertPiece(t, b, "a3", chess.Rook, chess.Light)
	assertNoPiece(t, b, "a1")

	assertParse(t, b, "Ra6")

	assertPiece(t, b, "a6", chess.Rook, chess.Dark)
	assertNoPiece(t, b, "a8")
}

func TestBoardMoveRookCapture(t *testing.T) {
	t.Parallel()

	b := chess.NewBoard()

	assertParse(t, b, "a4 e6 h4 e5 Ra3 e4 Rhh3 e3 Raxe3")

	assertPiece(t, b, "e3", chess.Rook, chess.Light)
	assertPiece(t, b, "h3", chess.Rook, chess.Light)
	assertNoPiece(t, b, "a3")
}

func TestBoardMoveRookInvalid(t *testing.T) {
	t.Parallel()

	b := chess.NewBoard()

	assertMoveError(t, b, "Rb2", "b2 blocked by white pawn")
	assertMoveError(t, b, "Rb1", "b1 blocked by white knight")
	assertMoveError(t, b, "Ra2", "a2 blocked by white pawn")
	assertMoveError(t, b, "Ra3", "no rook found that can move to a3")

	assertParse(t, b, "e3 e6 a4 d6 Ra3")

	// path blocked by pawn at d3
	assertMoveError(t, b, "Rh3", "no rook found that can move to h3")

	// ambiguous moves
	b = chess.NewBoard()

	assertParse(t, b, "a4 e6 h4 e5 Ra3 e4")

	assertMoveError(t, b, "Rh3", "move ambiguous: 2 rooks can move to h3")
	assertParse(t, b, "Rhh3")
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

	// path blocked by white pawn at d2
	assertMoveError(t, b, "Qd3", "no queen found that can move to d3")

	// TODO: ambiguous queen moves require pawn promotion
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

func TestBoardCheck(t *testing.T) {
	t.Parallel()

	b := chess.NewBoard()

	assert.False(t, b.InCheck())

	assertParse(t, b, "e4 e5 Qh5 Nc6 Qxf7")

	assert.True(t, b.InCheck())
	assert.True(t, strings.HasSuffix(b.Moves[len(b.Moves)-1], "+"), "check move should end with +")

	assertMoveError(t, b, "Nf6", "invalid move Nf6: king is in check")
	assertMoveError(t, b, "Ke7", "invalid move Ke7: king is in check")

	assertParse(t, b, "Kxf7")
}

func TestBoardPin(t *testing.T) {
	t.Parallel()

	b := chess.NewBoard()

	assertParse(t, b, "d4 e5 Nc3 Bb4")

	assertMoveError(t, b, "Ne4", "invalid move Ne4: king is in check")
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
