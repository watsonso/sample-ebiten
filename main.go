package main

import(
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/watsonso/sample-ebiten/2048"
)

func main() {
	game, err := twenty48.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowSize(twenty48.ScreenWidth, twenty48.ScreenHeight)
	ebiten.SetWindowTitle("2048(Ebiten Demo)")
	if err := ebiten.RunGame(gmae); err != nil {
		log.Fatal(err)
	}
}
