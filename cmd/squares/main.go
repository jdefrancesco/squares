package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jdefrancesco/squares/internal/game"
)

func main() {
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("Squares")

	if err := ebiten.RunGame(game.New()); err != nil {
		if err == ebiten.Termination {
			return
		}
		log.Fatal(err)
	}
}
