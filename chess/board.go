package chess

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"regexp"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type Board struct {
	tiles [8][8]*Piece
	turn  Color
	Moves []string
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
	if len(b.Moves) == 0 {
		return ""
	}

	var text string
	for i, m := range b.Moves {
		if i%2 == 0 {
			text += fmt.Sprintf("%d.%s", i/2+1, m)
		} else {
			text += fmt.Sprintf(" %s ", m)
		}
	}
	return fmt.Sprintf("`%s`", text)
}

func (b *Board) Parse(pgn string) error {
	var (
		moves = strings.Split(strings.Trim(pgn, " "), " ")
		re    = regexp.MustCompile(`[0-9]+\.`)
		err   error
	)

	for _, move := range moves {
		move = strings.Trim(move, " ")
		if move == "" {
			continue
		}

		// parse algebraic notation with numbers
		move = re.ReplaceAllString(move, "")

		if err = b.Move(move); err != nil {
			return err
		}
	}

	return nil
}

func (b *Board) Move(move string) error {
	var (
		to    string
		piece string
		// if the move is ambiguous, the originating square rank must be given
		// see https://en.wikipedia.org/wiki/Algebraic_notation_(chess)#Disambiguating_moves
		fromX          int
		fromY          int
		promotion      string
		castle         = false
		collisionPiece *Piece
		err            error
	)

	if parts := strings.Split(move, "="); len(parts) > 1 {
		promotion = parts[1]
		move = parts[0]
	}

	if piece, fromX, fromY, to, err = parseMove(move); err != nil {
		return err
	}

	if move == "O-O" || move == "O-O-O" {
		castle = true
		if b.turn == Dark {
			fromY = 0
			if move == "O-O" {
				to = "g8"
			} else {
				to = "c8"
			}
		}
	}

	// TODO: parse ambiguous captures for all pieces
	// TODO: parse checkmates e.g. e5#

	move_ := func() error {

		// collision detection
		if collisionPiece, err = b.getCollision(to); err != nil {
			return fmt.Errorf("invalid move %s: %v", move, err)
		} else if collisionPiece != nil {
			return fmt.Errorf("invalid move %s: position %s blocked by %s", move, to, collisionPiece)
		}

		switch strings.ToLower(piece) {
		case "p":
			return b.movePawn(to, fromX, fromY, promotion)
		case "r":
			return b.moveRook(to, false, fromX, fromY)
		case "b":
			return b.moveBishop(to, false)
		case "n":
			return b.moveKnight(to, fromX, fromY)
		case "q":
			return b.moveQueen(to, fromX, fromY)
		case "k":
			return b.moveKing(to, castle)
		default:
			return fmt.Errorf("invalid move %s: %v", move, err)
		}
	}

	if err = move_(); err != nil {
		return err
	}

	// is current player in check after move?
	if b.InCheck() {
		return fmt.Errorf("invalid move %s: king is in check", move)
	}

	if b.turn == Light {
		b.turn = Dark
	} else {
		b.turn = Light
	}

	// make sure the move is marked as a check if it was
	if b.InCheck() && !strings.HasSuffix(move, "+") && !strings.HasSuffix(move, "#") {
		move += "+"
	}

	b.Moves = append(b.Moves, move)

	return nil
}

func parseMove(move string) (string, int, int, string, error) {
	var (
		piece string
		fromX = -1
		fromY = -1
		to    string
		from  string
	)

	move = strings.TrimSuffix(move, "+")
	move = strings.TrimSuffix(move, "#")

	if move == "O-O" {
		return "K", 5, 7, "g1", nil
	}

	if move == "O-O-O" {
		return "K", 5, 7, "c1", nil
	}

	if strings.Contains(move, "x") {
		return parseCaptureMove(move)
	}

	if len(move) < 2 {
		return piece, fromX, fromY, to, fmt.Errorf("invalid move: %s", move)
	}

	if len(move) == 2 {
		// pawn move
		piece = "p"
		to = move
		return piece, fromX, fromY, to, nil
	}

	if len(move) == 3 {
		// unambiguous piece move
		piece = move[0:1]
		to = move[1:]
		return piece, fromX, fromY, to, nil
	}

	// last two characters are the target position
	to = move[len(move)-2:]

	// everything before is about origin
	from = move[:len(move)-2]

	// piece is first character
	piece = from[0:1]

	// we can ignore piece from here on
	from = from[1:]

	if len(from) == 1 {
		// rank or file given
		if from[0] >= 'a' && from[0] <= 'h' {
			// file given
			return piece, toFile(rune(from[0])), -1, to, nil
		}

		if from[0] >= '1' && from[0] <= '8' {
			// rank given
			return piece, -1, toRank(rune(from[0])), to, nil
		}

		return "", -1, -1, "", fmt.Errorf("invalid move: %s", move)
	}

	if len(from) == 2 {
		// file and rank given
		return piece, toFile(rune(from[0])), toRank(rune(from[1])), to, nil
	}

	return "", -1, -1, "", fmt.Errorf("invalid move: %s", move)
}

func parseCaptureMove(move string) (string, int, int, string, error) {
	var (
		parts = strings.Split(move, "x")
		piece string
		fromX = -1
		fromY = -1
		from  = parts[0]
		to    = parts[1]
	)

	if len(from) == 1 {
		// pawn move with rank given (exd4) or piece move (Nxe5)
		if strings.ToLower(from) == from {
			piece = "p"
			fromX = toFile(rune(from[0]))
			return piece, fromX, fromY, to, nil
		}

		piece = from
		return piece, fromX, fromY, to, nil
	}

	// piece is always first character
	piece = from[0:1]

	if len(from) == 2 {
		// rank or file given
		fromX = toFile(rune(from[1]))
		fromY = toRank(rune(from[1]))

		if fromX == -1 && fromY == -1 {
			return "", -1, -1, "", fmt.Errorf("invalid move: %s", move)
		}

		return piece, fromX, fromY, to, nil
	}

	if len(from) == 3 {
		// both file and rank given
		fromX = toFile(rune(from[1]))
		fromY = toRank(rune(from[1]))

		if fromX == -1 || fromY == -1 {
			return "", -1, -1, "", fmt.Errorf("invalid move: %s", move)
		}

		return piece, fromX, fromY, to, nil
	}

	return "", -1, -1, "", fmt.Errorf("invalid move: %s", move)
}

func toFile(r rune) int {
	if r >= 'a' && r <= 'h' {
		return int(r - 'a')
	}
	return -1
}

func toRank(r rune) int {
	if r >= '1' && r <= '8' {
		return int('8' - r)
	}
	return -1
}

func (b *Board) validateMove(search PieceName, fromX int, fromY int) func(*Piece, int, int) bool {
	return func(p *Piece, xPrev int, yPrev int) bool {
		// does piece originate from given origin?
		from := (fromX == -1 || xPrev == fromX) && (fromY == -1 || yPrev == fromY)

		// is this the piece we're looking for?
		found := p != nil && p.Name == search && p.Color == b.turn

		return from && found
	}
}

func (b *Board) InCheck() bool {
	var (
		kingX, kingY int
		p            *Piece
		x, y         int
	)

outerLoop:
	// find king
	for y = 0; y < 8; y++ {
		for x = 0; x < 8; x++ {
			p := b.getPiece(x, y)
			if p != nil && p.Name == King && p.Color == b.turn {
				kingX = x
				kingY = y
				break outerLoop
			}
		}
	}

	// check if any piece can attack king

	// ^
	x = kingX
	y = kingY
	for y = kingY - 1; y >= 0; y-- {
		p = b.getPiece(x, y)

		if p == nil {
			continue
		}

		if p.Color != b.turn && (p.Name == Rook || p.Name == Queen) {
			return true
		}

		// some other piece is blocking the way
		break
	}

	// ^>
	x = kingX
	y = kingY
	for x, y = kingX+1, kingY-1; x < 8 && y >= 0; x, y = x+1, y-1 {
		p = b.getPiece(x, y)

		if p == nil {
			continue
		}

		if p.Color != b.turn && (p.Name == Bishop || p.Name == Queen) {
			return true
		}

		break
	}

	// >
	x = kingX
	y = kingY
	for x = kingX + 1; x < 8; x++ {
		p = b.getPiece(x, y)

		if p == nil {
			continue
		}

		if p.Color != b.turn && (p.Name == Rook || p.Name == Queen) {
			return true
		}

		break
	}

	// v>
	x = kingX
	y = kingY
	for x, y = kingX+1, kingY+1; x < 8 && y < 8; x, y = x+1, y+1 {
		p = b.getPiece(x, y)

		if p == nil {
			continue
		}

		if p.Color != b.turn && (p.Name == Bishop || p.Name == Queen) {
			return true
		}

		break
	}

	// v
	x = kingX
	y = kingY
	for y = kingY + 1; y < 8; y++ {
		p = b.getPiece(x, y)

		if p == nil {
			continue
		}

		if p.Color != b.turn && (p.Name == Rook || p.Name == Queen) {
			return true
		}

		break
	}

	// <v
	x = kingX
	y = kingY
	for x, y = kingX-1, kingY+1; x >= 0 && y < 8; x, y = x-1, y+1 {
		p = b.getPiece(x, y)

		if p == nil {
			continue
		}

		if p.Color != b.turn && (p.Name == Bishop || p.Name == Queen) {
			return true
		}

		break
	}

	// <
	x = kingX
	y = kingY
	for x = kingX - 1; x >= 0; x-- {
		p = b.getPiece(x, y)

		if p == nil {
			continue
		}

		if p.Color != b.turn && (p.Name == Rook || p.Name == Queen) {
			return true
		}

		break
	}

	// <^
	x = kingX
	y = kingY
	for x, y = kingX-1, kingY-1; x >= 0 && y >= 0; x, y = x-1, y-1 {
		p = b.getPiece(x, y)

		if p == nil {
			continue
		}

		if p.Color != b.turn && (p.Name == Bishop || p.Name == Queen) {
			return true
		}

		break
	}

	// check for knights
	x = kingX + 1
	y = kingY - 2
	p = b.getPiece(x, y)
	if p != nil && p.Color != b.turn && p.Name == Knight {
		return true
	}

	x = kingX + 2
	y = kingY - 1
	p = b.getPiece(x, y)
	if p != nil && p.Color != b.turn && p.Name == Knight {
		return true
	}

	x = kingX + 2
	y = kingY + 1
	p = b.getPiece(x, y)
	if p != nil && p.Color != b.turn && p.Name == Knight {
		return true
	}

	x = kingX + 1
	y = kingY + 2
	p = b.getPiece(x, y)
	if p != nil && p.Color != b.turn && p.Name == Knight {
		return true
	}

	x = kingX - 1
	y = kingY + 2
	p = b.getPiece(x, y)
	if p != nil && p.Color != b.turn && p.Name == Knight {
		return true
	}

	x = kingX - 2
	y = kingY + 1
	p = b.getPiece(x, y)
	if p != nil && p.Color != b.turn && p.Name == Knight {
		return true
	}

	x = kingX - 2
	y = kingY - 1
	p = b.getPiece(x, y)
	if p != nil && p.Color != b.turn && p.Name == Knight {
		return true
	}

	x = kingX - 1
	y = kingY - 2
	p = b.getPiece(x, y)
	if p != nil && p.Color != b.turn && p.Name == Knight {
		return true
	}

	// check for pawns
	x = kingX - 1
	if b.turn == Light {
		y = kingY - 1
	} else {
		y = kingY + 1
	}
	p = b.getPiece(x, y)
	if p != nil && p.Color != b.turn && p.Name == Pawn {
		return true
	}

	x = kingX + 1
	if b.turn == Light {
		y = kingY - 1
	} else {
		y = kingY + 1
	}
	p = b.getPiece(x, y)
	if p != nil && p.Color != b.turn && p.Name == Pawn {
		return true
	}

	return false
}

func (b *Board) movePawn(position string, fromX int, fromY int, promotion string) error {
	var (
		toX   int
		toY   int
		yPrev int
		piece *Piece
		err   error
	)

	if toX, toY, err = getXY(position); err != nil {
		return err
	}

	if fromX != -1 {
		if b.turn == Light {
			fromY = toY + 1
		} else {
			fromY = toY - 1
		}

		piece = b.getPiece(fromX, fromY)
		if piece == nil || piece.Name != Pawn || piece.Color != b.turn {
			// not your pawn
			return fmt.Errorf("invalid capture move for pawn: %s", position)
		}

		if fromX != toX-1 && fromX != toX+1 || (b.turn == Light && fromY != toY+1) || (b.turn == Dark && fromY != toY-1) {
			// invalid capture move
			return fmt.Errorf("invalid capture move for pawn: %s", position)
		}

		b.tiles[fromX][fromY] = nil
		b.tiles[toX][toY] = piece
		if promotion != "" {
			return b.promotePawn(toX, toY, promotion)
		}
		return nil
	}

	// TODO: assert move is valid:
	//   * 2 moves from start position
	//   * 1 move otherwise
	//   * diagonal if attacking
	//   * no collision with other pieces

	if b.turn == Light {
		yPrev = toY + 1
	} else {
		yPrev = toY - 1
	}

	piece = b.tiles[toX][yPrev]
	if piece != nil && piece.Name == Pawn && piece.Color == b.turn {
		b.tiles[toX][yPrev] = nil
		b.tiles[toX][toY] = piece
		if promotion != "" {
			return b.promotePawn(toX, toY, promotion)
		}
		return nil
	}

	if b.turn == Light {
		yPrev = toY + 2
	} else {
		yPrev = toY - 2
	}

	piece = b.tiles[toX][yPrev]
	if piece != nil && piece.Name == Pawn && piece.Color == b.turn {
		b.tiles[toX][yPrev] = nil
		b.tiles[toX][toY] = piece
		if promotion != "" {
			return b.promotePawn(toX, toY, promotion)
		}
		return nil
	}

	return fmt.Errorf("no pawn found that can move to %s", position)
}

func (b *Board) promotePawn(x int, y int, name string) error {
	switch strings.ToLower(name) {
	case "q":
		b.tiles[x][y] = &Piece{Name: Queen, Color: b.turn}
	case "r":
		b.tiles[x][y] = &Piece{Name: Rook, Color: b.turn}
	case "b":
		b.tiles[x][y] = &Piece{Name: Bishop, Color: b.turn}
	case "n":
		b.tiles[x][y] = &Piece{Name: Knight, Color: b.turn}
	default:
		return fmt.Errorf("invalid promotion: %s", name)
	}

	return nil
}

func (b *Board) moveRook(position string, queen bool, fromX int, fromY int) error {
	var (
		x         int
		y         int
		xPrev     int
		yPrev     int
		p         *Piece
		validPrev []*Square
		err       error
	)

	if x, y, err = getXY(position); err != nil {
		return err
	}

	checkRookMove := b.validateMove(Rook, fromX, fromY)
	checkQueenMove := b.validateMove(Queen, fromX, fromY)
	checkMove := func(p *Piece, xPrev int, yPrev int) bool {
		return (!queen && checkRookMove(p, xPrev, yPrev)) || (queen && checkQueenMove(p, xPrev, yPrev))
	}

	xPrev = x
	yPrev = y
	for xPrev >= 0 && xPrev < 8 && yPrev >= 0 && yPrev < 8 {
		xPrev++
		p = b.getPiece(xPrev, yPrev)
		if p != nil {
			if checkMove(p, xPrev, yPrev) {
				validPrev = append(validPrev, &Square{X: xPrev, Y: yPrev})
			}
			break
		}
	}

	xPrev = x
	yPrev = y
	for xPrev >= 0 && xPrev < 8 && yPrev >= 0 && yPrev < 8 {
		yPrev--
		p = b.getPiece(xPrev, yPrev)
		if p != nil {
			if checkMove(p, xPrev, yPrev) {
				validPrev = append(validPrev, &Square{X: xPrev, Y: yPrev})
			} else {
				break
			}
		}
	}

	xPrev = x
	yPrev = y
	for xPrev >= 0 && xPrev < 8 && yPrev >= 0 && yPrev < 8 {
		xPrev--
		p = b.getPiece(xPrev, yPrev)
		if p != nil {
			if checkMove(p, xPrev, yPrev) {
				validPrev = append(validPrev, &Square{X: xPrev, Y: yPrev})
			} else {
				break
			}
		}
	}

	xPrev = x
	yPrev = y
	for xPrev >= 0 && xPrev < 8 && yPrev >= 0 && yPrev < 8 {
		yPrev++
		p = b.getPiece(xPrev, yPrev)
		if p != nil {
			if checkMove(p, xPrev, yPrev) {
				validPrev = append(validPrev, &Square{X: xPrev, Y: yPrev})
			} else {
				break
			}
		}
	}

	if len(validPrev) > 1 {
		return fmt.Errorf("move ambiguous: %d rooks can move to %s", len(validPrev), position)
	}

	if len(validPrev) == 1 {
		xPrev = validPrev[0].X
		yPrev = validPrev[0].Y
		p = b.getPiece(xPrev, yPrev)
		b.tiles[x][y] = p
		b.tiles[xPrev][yPrev] = nil
		return nil
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

func (b *Board) moveKnight(position string, fromX int, fromY int) error {
	var (
		x         int
		y         int
		xPrev     int
		yPrev     int
		p         *Piece
		validPrev []*Square
		err       error
	)

	if x, y, err = getXY(position); err != nil {
		return err
	}

	checkMove := b.validateMove(Knight, fromX, fromY)

	xPrev = x + 1
	yPrev = y - 2
	p = b.getPiece(xPrev, yPrev)
	if checkMove(p, xPrev, yPrev) {
		validPrev = append(validPrev, &Square{X: xPrev, Y: yPrev})
	}

	xPrev = x + 2
	yPrev = y - 1
	p = b.getPiece(xPrev, yPrev)
	if checkMove(p, xPrev, yPrev) {
		validPrev = append(validPrev, &Square{X: xPrev, Y: yPrev})
	}

	xPrev = x + 2
	yPrev = y + 1
	p = b.getPiece(xPrev, yPrev)
	if checkMove(p, xPrev, yPrev) {
		validPrev = append(validPrev, &Square{X: xPrev, Y: yPrev})
	}

	xPrev = x + 1
	yPrev = y + 2
	p = b.getPiece(xPrev, yPrev)
	if checkMove(p, xPrev, yPrev) {
		validPrev = append(validPrev, &Square{X: xPrev, Y: yPrev})
	}

	xPrev = x - 1
	yPrev = y + 2
	p = b.getPiece(xPrev, yPrev)
	if checkMove(p, xPrev, yPrev) {
		validPrev = append(validPrev, &Square{X: xPrev, Y: yPrev})
	}

	xPrev = x - 2
	yPrev = y + 1
	p = b.getPiece(xPrev, yPrev)
	if checkMove(p, xPrev, yPrev) {
		validPrev = append(validPrev, &Square{X: xPrev, Y: yPrev})
	}

	xPrev = x - 2
	yPrev = y - 1
	p = b.getPiece(xPrev, yPrev)
	if checkMove(p, xPrev, yPrev) {
		validPrev = append(validPrev, &Square{X: xPrev, Y: yPrev})
	}

	xPrev = x - 1
	yPrev = y - 2
	p = b.getPiece(xPrev, yPrev)
	if checkMove(p, xPrev, yPrev) {
		validPrev = append(validPrev, &Square{X: xPrev, Y: yPrev})
	}

	if len(validPrev) > 1 {
		return fmt.Errorf("move ambiguous: %d knights can move to %s", len(validPrev), position)
	}

	if len(validPrev) == 1 {
		xPrev = validPrev[0].X
		yPrev = validPrev[0].Y
		p = b.getPiece(xPrev, yPrev)
		b.tiles[x][y] = p
		b.tiles[xPrev][yPrev] = nil
		return nil
	}

	return fmt.Errorf("no knight found that can move to %s", position)
}

func (b *Board) moveQueen(position string, fromX int, fromY int) error {
	var (
		err error
	)

	// TODO: queen bishop moves might be ambiguous if there are multiple queens
	if err = b.moveBishop(position, true); err == nil {
		return nil
	}

	if err = b.moveRook(position, true, fromX, fromY); err == nil {
		return nil
	}

	return fmt.Errorf("no queen found that can move to %s", position)
}

func (b *Board) moveKing(position string, castle bool) error {
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

	if castle {
		// TODO: check if castle is allowed

		y := 7
		if b.turn == Dark {
			y = 0
		}

		if (b.turn == Light && position == "g1") || (b.turn == Dark && position == "g8") {
			// kingside castle

			king := b.getPiece(4, y)
			if king == nil || king.Color != b.turn || king.Name != King {
				return fmt.Errorf("invalid castle move")
			}

			if b.getPiece(5, y) != nil {
				return fmt.Errorf("invalid castle move")
			}

			if b.getPiece(6, y) != nil {
				return fmt.Errorf("invalid castle move")
			}

			rook := b.getPiece(7, y)
			if rook == nil || rook.Color != b.turn || rook.Name != Rook {
				return fmt.Errorf("invalid castle move")
			}

			b.tiles[6][y] = king
			b.tiles[4][y] = nil
			b.tiles[5][y] = rook
			b.tiles[7][y] = nil

			return nil
		}

		if (b.turn == Light && position == "c1") || (b.turn == Dark && position == "c8") {
			// queenside castle

			king := b.getPiece(4, y)
			if king == nil || king.Color != b.turn || king.Name != King {
				return fmt.Errorf("invalid castle move")
			}

			if b.getPiece(3, y) != nil {
				return fmt.Errorf("invalid castle move")
			}

			if b.getPiece(2, y) != nil {
				return fmt.Errorf("invalid castle move")
			}

			if b.getPiece(1, y) != nil {
				return fmt.Errorf("invalid castle move")
			}

			rook := b.getPiece(0, y)
			if rook == nil || rook.Color != b.turn || rook.Name != Rook {
				return fmt.Errorf("invalid castle move")
			}

			b.tiles[2][y] = king
			b.tiles[4][y] = nil
			b.tiles[3][y] = rook
			b.tiles[0][y] = nil

			return nil
		}
	}

	// ^
	xPrev = x + 0
	yPrev = y - 1
	piece = b.getPiece(xPrev, yPrev)
	if piece != nil && piece.Name == King && piece.Color == b.turn {
		b.tiles[xPrev][yPrev] = nil
		b.tiles[x][y] = piece
		return nil
	}

	// ^>
	xPrev = x + 1
	yPrev = y - 1
	piece = b.getPiece(xPrev, yPrev)
	if piece != nil && piece.Name == King && piece.Color == b.turn {
		b.tiles[xPrev][yPrev] = nil
		b.tiles[x][y] = piece
		return nil
	}

	// >
	xPrev = x + 1
	yPrev = y + 0
	piece = b.getPiece(xPrev, yPrev)
	if piece != nil && piece.Name == King && piece.Color == b.turn {
		b.tiles[xPrev][yPrev] = nil
		b.tiles[x][y] = piece
		return nil
	}

	// v>
	xPrev = x + 1
	yPrev = y + 1
	piece = b.getPiece(xPrev, yPrev)
	if piece != nil && piece.Name == King && piece.Color == b.turn {
		b.tiles[xPrev][yPrev] = nil
		b.tiles[x][y] = piece
		return nil
	}

	// v
	xPrev = x + 0
	yPrev = y + 1
	piece = b.getPiece(xPrev, yPrev)
	if piece != nil && piece.Name == King && piece.Color == b.turn {
		b.tiles[xPrev][yPrev] = nil
		b.tiles[x][y] = piece
		return nil
	}

	// <v
	xPrev = x - 1
	yPrev = y + 1
	piece = b.getPiece(xPrev, yPrev)
	if piece != nil && piece.Name == King && piece.Color == b.turn {
		b.tiles[xPrev][yPrev] = nil
		b.tiles[x][y] = piece
		return nil
	}

	// <
	xPrev = x - 1
	yPrev = y + 0
	piece = b.getPiece(xPrev, yPrev)
	if piece != nil && piece.Name == King && piece.Color == b.turn {
		b.tiles[xPrev][yPrev] = nil
		b.tiles[x][y] = piece
		return nil
	}

	// <^
	xPrev = x - 1
	yPrev = y - 1
	piece = b.getPiece(xPrev, yPrev)
	if piece != nil && piece.Name == King && piece.Color == b.turn {
		b.tiles[xPrev][yPrev] = nil
		b.tiles[x][y] = piece
		return nil
	}

	return fmt.Errorf("no king found that can move to %s", position)
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

	if len(runes) != 2 {
		return -1, -1, fmt.Errorf("square does not exist: %s", position)
	}

	posX = runes[0]
	posY = runes[1]

	if posX < 'a' || posX > 'h' {
		return -1, -1, fmt.Errorf("square does not exist: %s", position)
	}

	if posY < '1' || posY > '8' {
		return -1, -1, fmt.Errorf("square does not exist: %s", position)
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

func (b *Board) getCollision(position string) (*Piece, error) {
	var (
		x, y      int
		p         *Piece
		collision bool
		err       error
	)
	if x, y, err = getXY(position); err != nil {
		return nil, err
	}

	p = b.getPiece(x, y)

	// check if position is occupied by own piece
	collision = p != nil && p.Color == b.turn
	if collision {
		return p, nil
	}

	return nil, nil
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
