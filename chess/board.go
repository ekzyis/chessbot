package chess

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type Board struct {
	tiles [8][8]*Piece
	turn  Color
	moves []string
}

func NewBoard() *Board {
	board := &Board{turn: Light}

	board.mustSetPiece(Rook, Light, "a1")
	board.mustSetPiece(Knight, Light, "b1")
	board.mustSetPiece(Bishop, Light, "c1")
	board.mustSetPiece(Queen, Light, "d1")
	board.mustSetPiece(King, Light, "e1")
	board.mustSetPiece(Bishop, Light, "f1")
	board.mustSetPiece(Knight, Light, "g1")
	board.mustSetPiece(Rook, Light, "h1")
	board.mustSetPiece(Pawn, Light, "a2")
	board.mustSetPiece(Pawn, Light, "b2")
	board.mustSetPiece(Pawn, Light, "c2")
	board.mustSetPiece(Pawn, Light, "d2")
	board.mustSetPiece(Pawn, Light, "e2")
	board.mustSetPiece(Pawn, Light, "f2")
	board.mustSetPiece(Pawn, Light, "g2")
	board.mustSetPiece(Pawn, Light, "h2")

	board.mustSetPiece(Rook, Dark, "a8")
	board.mustSetPiece(Knight, Dark, "b8")
	board.mustSetPiece(Bishop, Dark, "c8")
	board.mustSetPiece(Queen, Dark, "d8")
	board.mustSetPiece(King, Dark, "e8")
	board.mustSetPiece(Bishop, Dark, "f8")
	board.mustSetPiece(Knight, Dark, "g8")
	board.mustSetPiece(Rook, Dark, "h8")
	board.mustSetPiece(Pawn, Dark, "a7")
	board.mustSetPiece(Pawn, Dark, "b7")
	board.mustSetPiece(Pawn, Dark, "c7")
	board.mustSetPiece(Pawn, Dark, "d7")
	board.mustSetPiece(Pawn, Dark, "e7")
	board.mustSetPiece(Pawn, Dark, "f7")
	board.mustSetPiece(Pawn, Dark, "g7")
	board.mustSetPiece(Pawn, Dark, "h7")

	return board
}

func NewGame(moves string) (*Board, error) {
	board := NewBoard()

	if err := board.Parse(moves); err != nil {
		return nil, err
	}

	return board, nil
}

func (b *Board) Save(filename string) error {
	var (
		file *os.File
		err  error
	)

	if file, err = os.Create(filename); err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, b.Image())
}

func (b *Board) Image() *image.RGBA {
	var (
		img   *image.RGBA
		piece *Piece
		bg    *image.Uniform
		rect  image.Rectangle
		p     = image.Point{0, 0}
	)

	img = image.NewRGBA(image.Rect(0, 0, 1024, 1024))

	for yi := 0; yi < 8; yi++ {
		for xi := 0; xi < 8; xi++ {
			x := xi * 128
			y := yi * 128
			rect = image.Rect(x, y, x+128, y+128)
			bg = image.NewUniform(getTileColor(xi, yi))
			draw.Draw(img, rect, bg, p, draw.Src)

			piece = b.tiles[xi][yi]
			if piece != nil {
				pieceImg := piece.Image
				if b.turn == Dark {
					pieceImg = flipImage(pieceImg)
				}
				draw.Draw(img, rect, pieceImg, p, draw.Over)
			}
		}
	}

	if b.turn == Dark {
		img = flipImage(img)
	}

	for yi := 0; yi < 8; yi++ {
		for xi := 0; xi < 8; xi++ {
			if b.turn == Light {
				drawCoordinate(img, xi, yi, false)
			}
			if b.turn == Dark {
				drawCoordinate(img, xi, yi, true)
			}
		}
	}

	return img
}

func drawCoordinate(img *image.RGBA, x, y int, flipped bool) {
	if x != 7 && y != 7 {
		return
	}

	var column, row string
	if y == 7 {
		switch x {
		case 0:
			column = "a"
		case 1:
			column = "b"
		case 2:
			column = "c"
		case 3:
			column = "d"
		case 4:
			column = "e"
		case 5:
			column = "f"
		case 6:
			column = "g"
		case 7:
			column = "h"
		}
	}

	if x == 7 {
		yRow := y
		if flipped {
			yRow = 7 - y
		}
		switch yRow {
		case 0:
			row = "8"
		case 1:
			row = "7"
		case 2:
			row = "6"
		case 3:
			row = "5"
		case 4:
			row = "4"
		case 5:
			row = "3"
		case 6:
			row = "2"
		case 7:
			row = "1"
		}
	}

	drawString := func(s string, origin fixed.Point26_6) {
		color := getTileColor(x, y)
		if !flipped && color == Light {
			color = Dark
		} else if !flipped {
			color = Light
		}
		// TODO: use SN font and make it bold
		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(color),
			Face: basicfont.Face7x13,
			Dot:  origin,
		}
		d.DrawString(s)
	}

	var origin fixed.Point26_6
	if column != "" {
		origin = fixed.P(x*128+5, (y+1)*128-5)
		drawString(column, origin)
	}

	if row != "" {
		origin = fixed.P((x+1)*128-12, y*128+15)
		drawString(row, origin)
	}
}

func (b *Board) SetPiece(name PieceName, color Color, position string) error {
	var (
		piece *Piece
		x     int
		y     int
		err   error
	)

	if len(position) != 2 {
		return fmt.Errorf("invalid position: %s", position)
	}

	if piece, err = NewPiece(name, color); err != nil {
		return err
	}

	if x, y, err = getXY(position); err != nil {
		return err
	}

	b.tiles[x][y] = piece

	return nil
}

func (b *Board) AlgebraicNotation() string {
	var text string
	for i, m := range b.moves {
		if i%2 == 0 {
			text += fmt.Sprintf("%d. %s", i/2+1, m)
		} else {
			text += fmt.Sprintf(" %s ", m)
		}
	}
	return fmt.Sprintf("`%s`", text)
}

func (b *Board) Parse(pgn string) error {
	var (
		moves = strings.Split(strings.Trim(pgn, " "), " ")
		err   error
	)

	for _, move := range moves {
		move = strings.Trim(move, " ")
		if move == "" {
			continue
		}
		if err = b.Move(move); err != nil {
			return err
		}
	}

	return nil
}

func (b *Board) Move(position string) error {
	var (
		err error
		// the column from which the piece is captured.
		// for example, this would be 'e' for exd4 and 'e' for Nexd4.
		captureFrom string
	)

	if strings.Contains(position, "x") {
		// capture move
		parts := strings.Split(position, "x")
		if len(parts) != 2 {
			return fmt.Errorf("invalid move: %s", position)
		}
		if len(parts[0]) == 2 {
			// example: Nexd4
			captureFrom = parts[0][1:]
		} else if len(parts[0]) == 1 {
			if strings.ToLower(parts[0]) == parts[0] {
				// example: exd4
				captureFrom = parts[0]
			} else {
				// example: Nxd4
				captureFrom = ""
			}
		} else {
			return fmt.Errorf("invalid move: %s", position)
		}

		position = strings.Replace(position, captureFrom+"x", "", 1)
	}

	// TODO: parse captures e.g. exd5 or Nxd5
	// TODO: parse promotions
	// TODO: parse checks e.g. e5+
	// TODO: parse checkmates e.g. e5#
	// TODO: parse O-O as kingside castle and O-O-O as queenside castle
	// TODO: make sure pinned pieces cannot move
	// TODO: make sure king is not in check after move
	//   ( this avoids moving into check and moving a piece that exposes the king to check e.g. pinned pieces )

	if len(position) == 2 {
		err = b.movePawn(position, captureFrom)
	} else if len(position) == 3 {
		switch strings.ToLower(position[0:1]) {
		case "r":
			err = b.moveRook(position[1:3], false)
		case "b":
			err = b.moveBishop(position[1:3], false)
		case "n":
			err = b.moveKnight(position[1:3])
		case "q":
			err = b.moveQueen(position[1:3])
		case "k":
			err = b.moveKing(position[1:3])
		default:
			err = fmt.Errorf("invalid move: %s", position)
		}
	} else {
		return fmt.Errorf("invalid move: %s", position)
	}

	if err != nil {
		return err
	}

	if b.turn == Light {
		b.turn = Dark
	} else {
		b.turn = Light
	}

	b.moves = append(b.moves, position)

	return nil
}

func (b *Board) movePawn(position string, captureFrom string) error {
	var (
		x     int
		y     int
		cX    int
		cY    int
		yPrev int
		piece *Piece
		err   error
	)

	if x, y, err = getXY(position); err != nil {
		return err
	}

	if captureFrom != "" {
		if cX, cY, err = getXY(captureFrom + position[0:1]); err != nil {
			return err
		}

		if b.turn == Light {
			cY = y + 1
		} else {
			cY = y - 1
		}

		piece = b.getPiece(cX, cY)
		if piece == nil || piece.Name != Pawn || piece.Color != b.turn {
			// not your pawn
			return fmt.Errorf("invalid capture move for pawn: %s", position)
		}

		if cX != x-1 && cX != x+1 || (b.turn == Light && cY != y+1) || (b.turn == Dark && cY != y-1) {
			// invalid capture move
			return fmt.Errorf("invalid capture move for pawn: %s", position)
		}

		b.tiles[cX][cY] = nil
		b.tiles[x][y] = piece
		return nil
	}

	// TODO: assert move is valid:
	//   * 2 moves from start position
	//   * 1 move otherwise
	//   * diagonal if attacking
	//   * no collision with other pieces

	if b.turn == Light {
		yPrev = y + 1
	} else {
		yPrev = y - 1
	}

	piece = b.tiles[x][yPrev]
	if piece != nil && piece.Name == Pawn && piece.Color == b.turn {
		b.tiles[x][yPrev] = nil
		b.tiles[x][y] = piece
		return nil
	}

	if b.turn == Light {
		yPrev = y + 2
	} else {
		yPrev = y - 2
	}

	piece = b.tiles[x][yPrev]
	if piece != nil && piece.Name == Pawn && piece.Color == b.turn {
		b.tiles[x][yPrev] = nil
		b.tiles[x][y] = piece
		return nil
	}

	return fmt.Errorf("no pawn found that can move to %s", position)
}

func (b *Board) moveRook(position string, queen bool) error {
	var (
		x     int
		y     int
		xPrev int
		yPrev int
		piece *Piece
		err   error
	)

	if x, y, err = getXY(position); err != nil {
		return err
	}

	xPrev = x
	yPrev = y
	for xPrev >= 0 && xPrev < 8 && yPrev >= 0 && yPrev < 8 {
		xPrev++
		piece = b.getPiece(xPrev, yPrev)
		if piece != nil {
			if ((!queen && piece.Name == Rook) || (queen && piece.Name == Queen)) && piece.Color == b.turn {
				b.tiles[xPrev][yPrev] = nil
				b.tiles[x][y] = piece
				return nil
			} else {
				// direction blocked by other piece
				break
			}
		}
	}

	xPrev = x
	yPrev = y
	for xPrev >= 0 && xPrev < 8 && yPrev >= 0 && yPrev < 8 {
		yPrev--
		piece = b.getPiece(xPrev, yPrev)
		if piece != nil {
			if ((!queen && piece.Name == Rook) || (queen && piece.Name == Queen)) && piece.Color == b.turn {
				b.tiles[xPrev][yPrev] = nil
				b.tiles[x][y] = piece
				return nil
			} else {
				break
			}
		}
	}

	xPrev = x
	yPrev = y
	for xPrev >= 0 && xPrev < 8 && yPrev >= 0 && yPrev < 8 {
		xPrev--
		piece = b.getPiece(xPrev, yPrev)
		if piece != nil {
			if ((!queen && piece.Name == Rook) || (queen && piece.Name == Queen)) && piece.Color == b.turn {
				b.tiles[xPrev][yPrev] = nil
				b.tiles[x][y] = piece
				return nil
			} else {
				break
			}
		}
	}

	xPrev = x
	yPrev = y
	for xPrev >= 0 && xPrev < 8 && yPrev >= 0 && yPrev < 8 {
		yPrev++
		piece = b.getPiece(xPrev, yPrev)
		if piece != nil {
			if ((!queen && piece.Name == Rook) || (queen && piece.Name == Queen)) && piece.Color == b.turn {
				b.tiles[xPrev][yPrev] = nil
				b.tiles[x][y] = piece
				return nil
			} else {
				break
			}
		}
	}

	return fmt.Errorf("no rook found that can move to %s", position)
}

func (b *Board) moveBishop(position string, queen bool) error {
	var (
		x     int
		y     int
		xPrev int
		yPrev int
		piece *Piece
		err   error
	)

	if x, y, err = getXY(position); err != nil {
		return err
	}

	xPrev = x
	yPrev = y
	for xPrev >= 0 && xPrev < 8 && yPrev >= 0 && yPrev < 8 {
		xPrev++
		yPrev--
		piece = b.getPiece(xPrev, yPrev)
		if piece != nil {
			if ((!queen && piece.Name == Bishop) || (queen && piece.Name == Queen)) && piece.Color == b.turn {
				b.tiles[xPrev][yPrev] = nil
				b.tiles[x][y] = piece
				return nil
			} else {
				// direction blocked by other piece
				break
			}
		}
	}

	xPrev = x
	yPrev = y
	for xPrev >= 0 && xPrev < 8 && yPrev >= 0 && yPrev < 8 {
		xPrev++
		yPrev++
		piece = b.getPiece(xPrev, yPrev)
		if piece != nil {
			if ((!queen && piece.Name == Bishop) || (queen && piece.Name == Queen)) && piece.Color == b.turn {
				b.tiles[xPrev][yPrev] = nil
				b.tiles[x][y] = piece
				return nil
			} else {
				break
			}
		}
	}

	xPrev = x
	yPrev = y
	for xPrev >= 0 && xPrev < 8 && yPrev >= 0 && yPrev < 8 {
		xPrev--
		yPrev++
		piece = b.getPiece(xPrev, yPrev)
		if piece != nil {
			if ((!queen && piece.Name == Bishop) || (queen && piece.Name == Queen)) && piece.Color == b.turn {
				b.tiles[xPrev][yPrev] = nil
				b.tiles[x][y] = piece
				return nil
			} else {
				break
			}
		}
	}

	xPrev = x
	yPrev = y
	for xPrev >= 0 && xPrev < 8 && yPrev >= 0 && yPrev < 8 {
		xPrev--
		yPrev--
		piece = b.getPiece(xPrev, yPrev)
		if piece != nil {
			if ((!queen && piece.Name == Bishop) || (queen && piece.Name == Queen)) && piece.Color == b.turn {
				b.tiles[xPrev][yPrev] = nil
				b.tiles[x][y] = piece
				return nil
			} else {
				break
			}
		}
	}

	return fmt.Errorf("no bishop found that can move to %s", position)
}

func (b *Board) moveKnight(position string) error {
	var (
		x     int
		y     int
		xPrev int
		yPrev int
		piece *Piece
		err   error
	)

	if x, y, err = getXY(position); err != nil {
		return err
	}

	xPrev = x + 1
	yPrev = y - 2
	piece = b.getPiece(xPrev, yPrev)
	if piece != nil && piece.Name == Knight && piece.Color == b.turn {
		b.tiles[xPrev][yPrev] = nil
		b.tiles[x][y] = piece
		return nil
	}

	xPrev = x + 2
	yPrev = y - 1
	piece = b.getPiece(xPrev, yPrev)
	if piece != nil && piece.Name == Knight && piece.Color == b.turn {
		b.tiles[xPrev][yPrev] = nil
		b.tiles[x][y] = piece
		return nil
	}

	xPrev = x + 2
	yPrev = y + 1
	piece = b.getPiece(xPrev, yPrev)
	if piece != nil && piece.Name == Knight && piece.Color == b.turn {
		b.tiles[xPrev][yPrev] = nil
		b.tiles[x][y] = piece
		return nil
	}

	xPrev = x + 1
	yPrev = y + 2
	piece = b.getPiece(xPrev, yPrev)
	if piece != nil && piece.Name == Knight && piece.Color == b.turn {
		b.tiles[xPrev][yPrev] = nil
		b.tiles[x][y] = piece
		return nil
	}

	xPrev = x - 1
	yPrev = y + 2
	piece = b.getPiece(xPrev, yPrev)
	if piece != nil && piece.Name == Knight && piece.Color == b.turn {
		b.tiles[xPrev][yPrev] = nil
		b.tiles[x][y] = piece
		return nil
	}

	xPrev = x - 2
	yPrev = y + 1
	piece = b.getPiece(xPrev, yPrev)
	if piece != nil && piece.Name == Knight && piece.Color == b.turn {
		b.tiles[xPrev][yPrev] = nil
		b.tiles[x][y] = piece
		return nil
	}

	xPrev = x - 2
	yPrev = y - 1
	piece = b.getPiece(xPrev, yPrev)
	if piece != nil && piece.Name == Knight && piece.Color == b.turn {
		b.tiles[xPrev][yPrev] = nil
		b.tiles[x][y] = piece
		return nil
	}

	xPrev = x - 1
	yPrev = y - 2
	piece = b.getPiece(xPrev, yPrev)
	if piece != nil && piece.Name == Knight && piece.Color == b.turn {
		b.tiles[xPrev][yPrev] = nil
		b.tiles[x][y] = piece
		return nil
	}

	return fmt.Errorf("no knight found that can move to %s", position)
}

func (b *Board) moveQueen(position string) error {
	var (
		err error
	)

	if err = b.moveBishop(position, true); err == nil {
		return nil
	}

	if err = b.moveRook(position, true); err == nil {
		return nil
	}

	return fmt.Errorf("no queen found that can move to %s", position)
}

func (b *Board) moveKing(position string) error {
	// TODO: implement king moves
	return nil
}

func (b *Board) mustSetPiece(name PieceName, color Color, position string) {
	if err := b.SetPiece(name, color, position); err != nil {
		log.Fatalf("cannot set piece %s: %v", name, err)
	}
}

func (b *Board) At(position string) *Piece {
	var (
		x   int
		y   int
		err error
	)
	if x, y, err = getXY(position); err != nil {
		return nil
	}
	return b.tiles[x][y]
}

func getXY(position string) (int, int, error) {
	var (
		posX rune
		posY rune
		x    int
		y    int
	)
	runes := []rune(position)
	posX = runes[0]
	posY = runes[1]

	if posX < 'a' || posX > 'h' {
		return -1, -1, fmt.Errorf("invalid move: %s", position)
	}

	if posY < '1' || posY > '8' {
		return -1, -1, fmt.Errorf("invalid move: %s", position)
	}

	// image origin (0,0) is at top-left corner (a8)
	x = int(posX - 'a')
	y = int('8' - posY)

	return x, y, nil
}

func (b *Board) getPiece(x int, y int) *Piece {
	if x < 0 || x >= 8 || y < 0 || y >= 8 {
		return nil
	}
	return b.tiles[x][y]
}

func getTileColor(x, y int) Color {
	if x%2 == y%2 {
		return Light
	} else {
		return Dark
	}
}

func flipImage(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	flipped := image.NewRGBA(bounds)

	// Flip the image vertically by reversing the Y coordinate
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			flipped.Set(x, bounds.Max.Y-y-1, img.At(x, y))
		}
	}

	return flipped
}
