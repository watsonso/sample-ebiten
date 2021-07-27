package twenty48

import (
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
)

var taskTerminated = errors.New("twenty48: task terminated")

type task func() error

// Board represents the game board.
type Board struct {
	size  int
	tiles map[*Tile]struct{}
	tasks []task
}

// NewBoard generates a new Board with giving a size.
func NewBoard(size int) (*Board, error) {
	b := &Board{
		size:  size,
		tiles: map[*Tile]struct{}{},
	}
	for i := 0; i < 2; i++ {
		if err := addRandomTile(b.tiles, b.size); err != nil {
			return nil, err
		}
	}
	return b, nil
}

func (b *Board) tileAt(x, y int) *Tile {
	return tileAt(b.tiles, x, y)
}

// Size returns the board size.
func (b *Board) Size() (int, int) {
	x := b.size*tileSize + (b.size+1)*tileMargin
	y := x
	return x, y
}

// Draw draws the board to the given boardImage.
func (b *Board) Draw(boardImage *ebiten.Image) {
	boardImage.Fill(frameColor)
	for j := 0; j < b.size; j++ {
		for i := 0; i < b.size; i++ {
			v := 0
			op := &ebiten.DrawImageOptions{}
			x := i*tileSize + (i+1)*tileMargin
			y := j*tileSize + (j+1)*tileMargin
			op.GeoM.Translate(float64(x), float64(y))
			r, g, b, a := colorToScale(tileBackgroundColor(v))
			op.ColorM.Scale(r, g, b, a)
			boardImage.DrawImage(tileImage, op)
		}
	}
	animatingTiles := map[*Tile]struct{}{}
	nonAnimatingTiles := map[*Tile]struct{}{}
	for t := range b.tiles {
		if t.IsMoving() {
			animatingTiles[t] = struct{}{}
		} else {
			nonAnimatingTiles[t] = struct{}{}
		}
	}
	for t := range nonAnimatingTiles {
		t.Draw(boardImage)
	}
	for t := range animatingTiles {
		t.Draw(boardImage)
	}
}
