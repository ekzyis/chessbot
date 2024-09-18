package chess

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"golang.org/x/image/draw"
)

type Piece struct {
	Name  PieceName
	Color color.Color
	Image image.Image
}

type PieceName string

const (
	Pawn   PieceName = "p"
	Knight PieceName = "n"
	Bishop PieceName = "b"
	Rook   PieceName = "r"
	Queen  PieceName = "q"
	King   PieceName = "k"
)

type Color color.Color

var (
	Light Color = color.RGBA{240, 217, 181, 255}
	Dark  Color = color.RGBA{181, 136, 99, 255}
)

func NewPiece(name PieceName, color Color) (*Piece, error) {
	var (
		colorSuffix string
		path        string
		file        *os.File
		img         image.Image
		dst         *image.RGBA
		err         error
	)

	colorSuffix = "l"
	if color == Light {
		colorSuffix = "l"
	} else if color == Dark {
		colorSuffix = "d"
	} else {
		return nil, fmt.Errorf("invalid color: %v", color)
	}

	path = fmt.Sprintf("assets/1024px-Chess_%s%st45.svg.png", name, colorSuffix)

	if file, err = os.Open(path); err != nil {
		return nil, err
	}
	defer file.Close()

	if img, err = png.Decode(file); err != nil {
		return nil, err
	}

	// source image for each piece is 1024x1024 and board is 8x8
	// so we need to scale each piece down to 128x128 (1024/8)
	dst = image.NewRGBA(image.Rect(0, 0, 128, 128))
	draw.CatmullRom.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)

	return &Piece{Name: name, Color: color, Image: dst}, nil
}
