package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jdefrancesco/squares/internal/game"
)

// main configures the game window and runs the Ebiten loop, treating a requested
// quit as normal termination and reporting other errors.
func main() {
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("Squares")

	g := game.New()
	err := ebiten.RunGame(g)
	if closeErr := g.Close(); closeErr != nil {
		log.Printf("Could not save personal best: %v", closeErr)
	}
	if err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
