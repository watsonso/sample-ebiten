package twenty48

import (
	"errors"
	"image/color"
	"log"
	"math/rand"
	"strconv"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

var (
	mplusSmallFont  font.Face
	mplusNormalFont font.Face
	mplusBigFont    font.Face
)

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusSmallFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	mplusBigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

// TileData represents a tile information like a value and a position.
type TileData struct {
	value int
	x     int
	y     int
}

// Tile represents s tile information including TileData and animation states.
type Tile struct {
	current TileData

	// next represents a next tile information after moving.
	// next is empty when the tile is not about to move.
	next TileData

	movingCount       int
	startPoppingCount int
	poppingCount      int
}

// NewTile creates a new Tile object.
func NewTile(value int, x, y int) *Tile {
	return &Tile{
		current: TileData{
			value: value,
			x:     x,
			y:     y,
		},
		startPoppingCount: maxPoppingCount,
	}
}

// IsMoving returns a boolean value indicating if the tile is animating.
func (t *Tile) IsMoving() bool {
	return 0 < t.movingCount
}

func tileAt(tiles map[*Tile]struct{}, x, y int) *Tile {
	var result *Tile
	for t := range tiles {
		if t.current.x != x || t.current.y != y {
			continue
		}
		if result != nil {
			panic("not reach")
		}
		result = t
	}
	return result
}

const (
	maxMovingCount  = 5
	maxPoppingCount = 6
)

func addRandomTile(tiles map[*Tile]struct{}, size int) error {
	cells := make([]bool, size*size)
	for t := range tiles {
		if t.IsMoving() {
			panic("not reach")
		}
		i := t.current.x + t.current.y*size
		cells[i] = true
	}
	availableCells := []int{}
	for i, b := range cells {
		if b {
			continue
		}
		availableCells = append(availableCells, i)
	}
	if len(availableCells) == 0 {
		return errors.New("twenty48: there is no space to add a new tile")
	}
	c := availableCells[rand.Intn(len(availableCells))]
	v := 2
	if rand.Intn(10) == 0 {
		v = 4
	}
	x := c % size
	y := c / size
	t := NewTile(v, x, y)
	tiles[t] = struct{}{}
	return nil
}

func colorToScale(clr color.Color) (float64, float64, float64, float64) {
	r, g, b, a := clr.RGBA()
	rf := float64(r) / 0xffff
	gf := float64(g) / 0xffff
	bf := float64(b) / 0xffff
	af := float64(a) / 0xffff
	// Convert to non-premultiplied alpha components.
	if 0 < af {
		rf /= af
		gf /= af
		bf /= af
	}
	return rf, gf, bf, af
}

func mean(a, b int, rate float64) int {
	return int(float64(a)*(1-rate) + float64(b)*rate)
}

func meanF(a, b float64, rate float64) float64 {
	return a*(1-rate) + b*rate
}

const (
	tileSize   = 80
	tileMargin = 4
)

var (
	tileImage = ebiten.NewImage(tileSize, tileSize)
)

func init() {
	tileImage.Fill(color.White)
}

// Draw draws the current tile to the given boardImage.
func (t *Tile) Draw(boardImage *ebiten.Image) {
	i, j := t.current.x, t.current.y
	ni, nj := t.next.x, t.next.y
	v := t.current.value
	if v == 0 {
		return
	}
	op := &ebiten.DrawImageOptions{}
	x := i*tileSize + (i+1)*tileMargin
	y := j*tileSize + (j+1)*tileMargin
	nx := ni*tileSize + (ni+1)*tileMargin
	ny := nj*tileSize + (nj+1)*tileMargin
	switch {
	case 0 < t.movingCount:
		rate := 1 - float64(t.movingCount)/maxMovingCount
		x = mean(x, nx, rate)
		y = mean(y, ny, rate)
	case 0 < t.startPoppingCount:
		rate := 1 - float64(t.startPoppingCount)/float64(maxPoppingCount)
		scale := meanF(0.0, 1.0, rate)
		op.GeoM.Translate(float64(-tileSize/2), float64(-tileSize/2))
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(float64(tileSize/2), float64(tileSize/2))
	case 0 < t.poppingCount:
		const maxScale = 1.2
		rate := 0.0
		if maxPoppingCount*2/3 <= t.poppingCount {
			// 0 to 1
			rate = 1 - float64(t.poppingCount-2*maxPoppingCount/3)/float64(maxPoppingCount/3)
		} else {
			// 1 to 0
			rate = float64(t.poppingCount) / float64(maxPoppingCount*2/3)
		}
		scale := meanF(1.0, maxScale, rate)
		op.GeoM.Translate(float64(-tileSize/2), float64(-tileSize/2))
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(float64(tileSize/2), float64(tileSize/2))
	}
	op.GeoM.Translate(float64(x), float64(y))
	r, g, b, a := colorToScale(tileBackgroundColor(v))
	op.ColorM.Scale(r, g, b, a)
	boardImage.DrawImage(tileImage, op)
	str := strconv.Itoa(v)

	f := mplusBigFont
	switch {
	case 3 < len(str):
		f = mplusSmallFont
	case 2 < len(str):
		f = mplusNormalFont
	}

	bound, _ := font.BoundString(f, str)
	w := (bound.Max.X - bound.Min.X).Ceil()
	h := (bound.Max.Y - bound.Min.Y).Ceil()
	x = x + (tileSize-w)/2
	y = y + (tileSize-h)/2 + h
	text.Draw(boardImage, str, f, x, y, tileColor(v))
}
