package twenty48

// "github.com/hajimehoshi/ebiten/v2"
// "github.com/hajimehoshi/ebiten/v2/inpututil"

// Dir represents a direction.
type Dir int

// Input represents the current key states.
type Input struct {
	mouseDir Dir
}

// NewInput generates a new Input object.
func NewInput() *Input {
	return &Input{}
}
