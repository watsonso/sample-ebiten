package twenty48

import (
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
)

var taskTerminated = errors.New("twenty48: task terminated")

type task func() error

// Board represents the game board.
type Board struct {
	size int
	// tiles map[*Tile]struct{}
	// tasks []task
}

// NewBoard generates a new Board with giving a size.
func NewBoard(size int) (*Board, error) {
	b := &Board{
		size: size,
		// tiles: map[*Tile]struct{}{},
	}
	// for i := 0; i < 2; i++ {
	// 	if err := addRandomTile(b.tiles, b.size); err != nil {
	// 		return nil, err
	// 	}
	// }
	return b, nil
}

func (b *Board) Size() (int, int) {
	w := 1
	h := 1
	return w, h
}

func (b *Board) Draw(boardImage *ebiten.Image){
	boardImage = nil
}
